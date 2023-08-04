package ada

import (
	"ada-telegram-bot/pkg/models"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Основная структура бота
type AdaBot struct {
	bot      *tgbotapi.BotAPI
	db       models.DB
	sessions map[int64]*models.Session
}

// Создание телеграмм бота.
func NewAdaBot(db models.DB) (*AdaBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("NewBotTelegram: tgbotapi.NewBotAPI: %w", err)
	}
	bot.Debug = false

	tgBot := AdaBot{
		bot:      bot,
		db:       db,
		sessions: make(map[int64]*models.Session, 10),
	}

	return &tgBot, nil
}

// Запуск бота
func (a *AdaBot) Run() error {
	// Запуск алерт менеджера
	go a.alertTicker()

	//TODO: Режим получения хуков (выбор)
	return a.runUpdater()
}

// Запуск апдейтора сообщений
func (a *AdaBot) runUpdater() error {
	log.Printf("Authorized on account %s", a.bot.Self.UserName)

	// Инициализация канала событий
	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 30

	if err := a.handler(a.bot.GetUpdatesChan(updater)); err != nil {
		return err
	}
	return nil
}

// Обработчики сообщений
func (a *AdaBot) handler(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		userId := update.Message.Chat.ID
		var msg *tgbotapi.Message
		var data string

		// Получение данных
		switch {
		// Обработка сообщений
		case update.Message != nil:
			msg = update.Message
			data = msg.Text
		// Обработка команд
		case update.Message != nil && update.Message.IsCommand():
			msg = update.Message
			data = msg.Command()
		// Обработка CallbackQuery
		case update.CallbackQuery != nil:
			msg = update.CallbackQuery.Message
			data = update.CallbackQuery.Data
		}

		// Логирование сообщения
		log.Printf("\nuserId: %d; messageId: %d, data: %s;\n", userId, msg.MessageID, data)

		// Проврека пользователя
		if err := a.createUser(msg); err != nil {
			log.Printf("user_id: %d; error a.createUser: %v", userId, err)
			continue
		}

		// Добавление сообщения пользователя в БД.
		if err := a.db.AddUserMessageId(userId, msg.MessageID, "user"); err != nil {
			log.Printf("user_id: %d; error db.AddUserMessageId: %v", userId, err)
			continue
		}

		// Обновление даты последней активности.
		if err := a.db.UpdateLastActive(userId); err != nil {
			log.Printf("user_id: %d; error db.UpdateLastActive: %v", userId, err)
			continue
		}

		// Отправка данных в обработчик
		go a.handlerSession(userId, msg.MessageID, data)

	}

	return fmt.Errorf("updates channel closed")
}

// Обработчик сессий.
// Обрабатывает полученные данные в зависимости от состояния сессии
func (a *AdaBot) handlerSession(userId int64, messageId int, data string) {
	session, err := a.db.GetLastSession(userId)
	if err != nil {
		log.Printf("userId: %d; messageId: %d, error a.getSession: %v", userId, messageId, err)
		return
	}

	// Стоп слово (сброс сессии)
	if data == "start" || data == "/start" {
		if err := a.start(userId); err != nil {
			log.Printf("userId: %d; messageId: %d, error a.start: %v", userId, messageId, err)
			return
		}
		return
	}

	switch {
	// Обработка функции
	case session.F != nil:
		if err := a.handlerFunc(userId, data); err != nil {
			log.Printf("userId: %d; messageId: %d, error db.handlerFunc: %v", userId, messageId, err)
			return
		}
	// Обработка цепочек
	case session.C != nil:
		if err := a.handlerChain(userId, data); err != nil {
			log.Printf("userId: %d; messageId: %d, error db.handlerChain: %v", userId, messageId, err)
			return
		}
	// Обработка меню
	case session.M != nil:
		if err := a.handlerMenu(userId, data); err != nil {
			log.Printf("userId: %d; messageId: %d, error db.handlerMenu: %v", userId, messageId, err)
			return
		}
	default:
		if err := a.sendErrorMessage(userId); err != nil {
			log.Printf("user_id: %d; error sendErrorMessage: %v", userId, err)
			return
		}
	}
}

// Отправка сообщения пользователю.
// Сохранение id сообщения в БД.
func (a *AdaBot) sendMessage(userId int64, c tgbotapi.Chattable) error {
	botMsg, err := a.bot.Send(c)
	if err != nil {
		return err
	}

	if err := a.db.AddUserMessageId(userId, botMsg.MessageID, "bot"); err != nil {
		return err
	}

	return nil
}

// Проврека и регистрация пользователя в бд
func (a *AdaBot) createUser(msg *tgbotapi.Message) error {
	// Регистрация пользователя если его нет.
	if err := a.db.CreateUser(msg.Chat.ID, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return fmt.Errorf("error db.CreateUser: %w", err)
	}
	return nil
}

// Сброс сессии.
func (a *AdaBot) start(userId int64) error {
	// Создание стартовой сессии
	if err := a.resetSession(userId); err != nil {
		return fmt.Errorf("error resetSession: %w", err)
	}

	// Очистка чата.
	if err := a.cleareChat(userId); err != nil {
		return err
	}

	// Отправка стартовых сообщений.

	// TODO: Отправка adMsg.
	// if viper.GetBool("ada_bot.ad_message") {
	// 	if err := a.sendAdMessage(userId); err != nil {
	// 		return err
	// 	}
	// }

	// TODO Отправка infoMsg.

	// Отправка startMsg.
	if err := a.sendStartMessage(userId); err != nil {
		return err
	}

	return nil
}

// Отправка startMessage
func (a *AdaBot) sendStartMessage(userId int64) error {
	text := `📓 <b>Возможности бота:</b>`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Рекламные интеграции", "ad"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Статистика", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Информация", "info"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тех. поддержка", "help"),
		),
	)
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.ReplyMarkup = keyboard

	if err := a.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error send startMessage: %w", err)
	}

	return nil
}

// Отправка adMessage
// TODO: не реализован
func (a *AdaBot) sendAdMessage(userId int64) error {
	// Создание botMsg adMessage.
	text := `📓 <b>💵 РЕКЛАМА </b>`
	// keyboard := tgbotapi.NewInlineKeyboardMarkup(
	// 	tgbotapi.NewInlineKeyboardRow(
	// 		tgbotapi.NewInlineKeyboardButtonData("Управление событиями", "ad_event"),
	// 	),
	// )
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	// botMsg.ReplyMarkup = keyboard

	// Отправка botMsg adMessage.
	newAdMessage, err := a.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new adMessage: %w", err)
	}

	// Сохранение adMessageId.
	if err := a.db.AddUserMessageId(userId, newAdMessage.MessageID, "ad"); err != nil {
		return err
	}

	// Удаление если возможно старого startMessage.
	adMessageId, err := a.db.GetAdMessageId(userId)
	if err != nil {
		log.Println("b.db.GetStartmessageId startMenu error: ", err)
	}
	a.cleareMessage(userId, adMessageId)

	// Установка нового adMessage.
	if err := a.db.UpdateAdMessageId(userId, newAdMessage.MessageID); err != nil {
		return err
	}

	return nil
}

// Отправка в чат сообщения о повторной попытке
func (a *AdaBot) sendErrorMessage(userId int64) error {
	// Отправка сообщения
	botMsg := tgbotapi.NewMessage(userId, "К сожалению что то пошло не так 🥲. Попробуйте повторно <b>/start</b> ")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard
	if err := a.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error send errorMessage: %w", err)
	}

	// Получение сессии
	session, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("error db.GetLastSession: %w", err)
	}

	// Логирование сессии
	log.Printf("user_id: %d, session: %v", userId, session)
	return nil
}

// Очистка сообщения
func (a *AdaBot) cleareMessage(userId int64, messageId int) error {
	if err := a.db.DeleteUsermessageId(messageId); err != nil {
		return err
	}

	deleteMsg := tgbotapi.NewDeleteMessage(userId, messageId)
	if _, err := a.bot.Send(deleteMsg); err != nil {
		return fmt.Errorf("error cleare messageId%d: %w", messageId, err)
	}
	return nil
}

// Очистка чата
func (a *AdaBot) cleareChat(userId int64) error {
	// Получение всех messageId
	messageIds, err := a.db.GetUserMessageIds(userId)
	if err != nil {
		return err
	}

	// Удаление всех сообщений кроме последнего
	for _, messageId := range messageIds {
		a.cleareMessage(userId, messageId)
	}

	return nil
}

// Откат действия
// TODO не реализован
func (a *AdaBot) backStep(userId int64) error {
	// Удаление крайней сессии
	if err := a.db.DeleteLastSession(userId); err != nil {
		return fmt.Errorf("error db.DeleteLastSession: %w", err)
	}
	//

	return nil
}

// Сброс сессии.
func (a *AdaBot) resetSession(userId int64) error {
	// Очистка сессий
	if err := a.db.DeleteAllSession(userId); err != nil {
		return fmt.Errorf("error db.DeleteAllSession: %w", err)
	}

	// Создание стартовой сессии
	startSession := models.Session{
		M:    []models.Menu{"start"},
		Data: make(map[string]interface{}),
	}

	// Сохранение в локальные сессии
	a.sessions[userId] = &startSession

	if err := a.saveSession(userId); err != nil {
		
		return fmt.Errorf("error saveSession: %w", err)
	}

	return nil
}

// Сохраняет текущую сессию в БД
func (a *AdaBot) saveSession(userId int64) error {
	if _, err := a.db.SaveSassion(userId, a.sessions[userId]); err != nil {
		return fmt.Errorf("error db.AddSession: %w", err)
	}

	return nil
}

// Получить данные по ключу
func (a *AdaBot) getDataSession(userId int64, key string) (value interface{}) {
	// Получение сессии
	session, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("error a.getSession: %w", err)
	}

	return session.Data[key]
}

// Добавить данные по ключу
func (a *AdaBot) setDataSession(userId int64, key string, value interface{}) error {
	// Получение сессии
	session, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("error a.getSession: %w", err)
	}

	// Обновление сессии
	session.Data[key] = value

	// Сохранение сессии
	_, err = a.db.SaveSassion(userId, session)
	if err != nil {
		return fmt.Errorf("error db.AddSession: %w", err)
	}

	return nil
}

// Открытие нового меню
func (a *AdaBot) addMenu(userId int64, m models.Menu) error {
	// Получение сессии
	session, err := a.db.GetLastSession(userId)
	if err != nil {
		return fmt.Errorf("error a.getSession: %w", err)
	}

	// Обновление сессии
	session.M = append(session.M, m)

	// Сохранение сессии в бд
	_, err = a.db.SaveSassion(userId, session)
	if err != nil {
		return fmt.Errorf("error db.AddSession: %w", err)
	}

	return nil
}
