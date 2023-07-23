package telegram

import (
	"AdaTelegramBot/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура телеграмм бота.
type BotTelegram struct {
	bot                  *tgbotapi.BotAPI
	db                   models.TelegramBotDB
	sessions             map[int64]*models.Session
	adEventsCache        map[int64][][]models.AdEvent // Хэш-таблица полученных из БД событий.
	adEventCreatingCache map[int64]*models.AdEvent    // Хэш-таблица создаваемых ad событий.
}

// Создание телеграмм бота.
func NewBotTelegram(db models.TelegramBotDB) (*BotTelegram, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("NewBotTelegram: tgbotapi.NewBotAPI: %w", err)
	}
	bot.Debug = false

	tgBot := BotTelegram{
		bot:                  bot,
		db:                   db,
		sessions:             make(map[int64]*models.Session),
		adEventsCache:        make(map[int64][][]models.AdEvent),
		adEventCreatingCache: make(map[int64]*models.AdEvent),
	}

	return &tgBot, nil
}

// Инициализация канала событий.
func (b *BotTelegram) InitUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	return b.bot.GetUpdatesChan(u)
}

// Обработчики сообщений.
func (b *BotTelegram) handlerUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		// Обработка команд.
		if update.Message != nil && update.Message.IsCommand() {
			// Добавление сообщения пользователя в БД.
			if err := b.db.AddUserMessageId(update.Message.Chat.ID,
				update.Message.MessageID); err != nil {
				log.Println("critical error: error AddUserMessageId to data base: ", err)
				return err
			}

			go b.handlerCommand(update.Message)
			continue
		}

		// Обработка сообщений.
		if update.Message != nil {
			// Добавление сообщения пользователя в БД.
			if err := b.db.AddUserMessageId(update.Message.Chat.ID,
				update.Message.MessageID); err != nil {
				log.Println("critical error: error AddUserMessageId to data base: ", err)
				return err
			}

			go b.handlerMessage(update.Message)
			continue
		}

		// Обработка CallbackQuery.
		if update.CallbackQuery != nil {
			// Добавление сообщения пользователя в БД.
			if err := b.db.AddUserMessageId(update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID); err != nil {
				log.Println("critical error: error AddUserMessageId to data base: ", err)
				return err
			}

			go b.handlerCbq(update.CallbackQuery)
			continue
		}
	}

	return fmt.Errorf("updates channel closed")
}

// Запуск апдейтера.
func (b *BotTelegram) StartBotUpdater() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates := b.InitUpdatesChannel()
	go b.alertTicker()
	if err := b.handlerUpdates(updates); err != nil {
		return err
	}
	return nil
}

// Получение хэша ad события.
func (b *BotTelegram) getAdEventCreatingCache(userId int64) (*models.AdEvent, error) {
	adEvent, ok := b.adEventCreatingCache[userId]
	if ok {
		return adEvent, nil
	}

	if err := b.sendRequestRestartMsg(userId); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("adEvent cache not found")
}

// Получение хэша ad событий.
func (b *BotTelegram) getAdEventsCache(userId int64) ([][]models.AdEvent, error) {
	adEvents, ok := b.adEventsCache[userId]
	if ok {
		return adEvents, nil
	}

	if err := b.sendRequestRestartMsg(userId); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("adEvents cache not found")
}

// Отправка в чат сообщения о повторной попытке.
func (b *BotTelegram) sendRequestRestartMsg(userId int64) error {
	b.db.SetStepUser(userId, "start")
	botMsg := tgbotapi.NewMessage(userId, "К сожалению что то пошло не так 🥲. Попробуйте повторно <b>/start</b> ")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error send message in sendRestartMessage: %w", err)
	}
	return nil
}

// Очистка сообщения.
func (b *BotTelegram) cleareMessage(userId int64, messageId int) error {
	if err := b.db.DeleteUsermessageId(messageId); err != nil {
		return err
	}

	deleteMsg := tgbotapi.NewDeleteMessage(userId, messageId)
	if _, err := b.bot.Send(deleteMsg); err != nil {
		return fmt.Errorf("error cleare messageId%d: %w", messageId, err)
	}
	return nil
}

// Очистка чата.
func (b *BotTelegram) cleareAllChat(userId int64) error {
	startMessageId, err := b.db.GetStartMessageId(userId)
	if err != nil {
		return err
	}

	adMessageId, err := b.db.GetAdMessageId(userId)
	if err != nil {
		return err
	}

	infoMessageId, err := b.db.GetAdMessageId(userId)
	if err != nil {
		return err
	}

	// Получение всех messageId.
	messageIds, err := b.db.GetUserMessageIds(userId)
	if err != nil {
		return err
	}

	// Удаление всех сообщений кроме startMessage / adMessage / infoMessage.
	for _, messageId := range messageIds {
		if messageId == startMessageId || messageId == adMessageId || messageId == infoMessageId {
			continue
		}
		b.cleareMessage(userId, messageId)
	}

	return nil
}

// Отправка сообщения пользователю.
func (b *BotTelegram) sendMessage(userId int64, c tgbotapi.Chattable) error {
	botMsg, err := b.bot.Send(c)
	if err != nil {
		return err
	}

	if err := b.db.AddUserMessageId(userId, botMsg.MessageID); err != nil {
		return err
	}

	return nil
}

// Отправка оповещения пользователю.
func (b *BotTelegram) sendAlertMessage(userId int64, c tgbotapi.Chattable) error {
	botMsg, err := b.bot.Send(c)
	if err != nil {
		return err
	}

	// Добавления ID сообщения в бд.
	if err := b.db.AddUserMessageId(userId, botMsg.MessageID); err != nil {
		return err
	}

	// Обновление даты оповещения.
	if err := b.db.UpdateTimeLastAlert(userId, time.Now()); err != nil {
		return err
	}

	return nil
}

// Иницализация сессии
func (b *BotTelegram) initSessions(userId int64) {
	b.sessions[userId] = &models.Session{
		Cache: make(map[string]interface{}),
	}
}

// Если ad событие полностью заполенно - возвращается true. Иначе false.
func fullDataAdEvent(ae *models.AdEvent) bool {
	if ae.UserId == 0 {
		log.Println("not found ae.UserId event")
		return false
	}

	if ae.Type == "" {
		log.Println("not found ae.Type event")
		return false
	}

	if ae.CreatedAt == "" {
		log.Println("not found ae.CreatedAt event")
		return false
	}

	if ae.DateStart == "" {
		log.Println("not found ae.DateStart event")
		return false
	}

	return true
}

// Получение данных из кэша.
func (b *BotTelegram) fromCache(userId int64, key string) (interface{}, bool) {
	if b.sessions == nil {
		return nil, false
	}

	session, ok := b.sessions[userId]
	if !ok {
		return nil, false
	}

	value, ok := session.Cache[key]
	if !ok {
		return nil, false
	}

	return value, true
}

// Запись в кэш.
func (b *BotTelegram) toCache(userId int64, key string, value interface{}) bool {
	if b.sessions == nil {
		return false
	}

	session, ok := b.sessions[userId]
	if !ok {
		return false
	}

	session.Cache[key] = value

	return true
}
