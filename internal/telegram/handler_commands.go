package telegram

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// Обработчик команд.
func (b *BotTelegram) handlerCommand(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	fmt.Printf("Info %s; user=%s; CMD=%s;\n", time.Now().Format("2006-01-02 15:04:05.999"), msg.From.UserName, msg.Command())
	switch msg.Command() {
	case "start":
		if err := b.cmdStart(msg); err != nil {
			log.Println("error in cmdStart: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
		return nil
	default:
		if err := b.handlerMessage(msg); err != nil {
			log.Println("error in handlerMessage: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
		// botMsg := tgbotapi.NewMessage(userId, `Неизвестная команда 🥲`)
		// botMsg.ParseMode = tgbotapi.ModeHTML
		// if err := b.sendMessage(userId, botMsg); err != nil {
		// 	return fmt.Errorf("error send unknow command error: %w", err)
		// }
		return nil
	}
}

// Команда /start
func (b *BotTelegram) cmdStart(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	b.initSessions(userId)

	// Регистрация пользователя если его нет.
	if err := b.db.DefaultUserCreation(userId, msg.Chat.UserName, msg.Chat.FirstName); err != nil {
		return err
	}

	// Очистка кэша пользователя.
	if err := b.clearCacheOfUser(userId); err != nil {
		return err
	}

	// Обновление даты последней активности.
	if err := b.db.UpdateLastActive(userId); err != nil {
		return err
	}

	// Отправка adMsg.
	if viper.GetBool("ada_bot.ad_message") {
		if err := b.sendAdMessage(userId); err != nil {
			return err
		}
	} else {
		if err := b.db.UpdateAdMessageId(userId, 0); err != nil {
			return err
		}
	}

	// TODO Отправка infoMsg.

	// Отправка startMsg.
	if err := sendStartMessage(b, userId); err != nil {
		return err
	}

	// Очистка чата.
	if err := b.cleareAllChat(userId); err != nil {
		return err
	}

	return nil
}

// Отправка startMessage.
func sendStartMessage(b *BotTelegram, userId int64) error {
	// Установка шага пользователя.
	if err := b.db.SetStepUser(userId, "start"); err != nil {
		return err
	}

	// Создание botMsg startMessage.
	text := `📓 <b>Возможности телеграмм бота:</b>`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать событие", "ad_event.create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Просмотреть события", "ad_event.view"),
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

	newStartMessage, err := b.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new startMessage: %w", err)
	}

	// Сохранение startMessageId.
	if err := b.db.AddUserMessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	// Удаление если возможно старого startMessage.
	startMessageId, err := b.db.GetStartMessageId(userId)
	if err != nil {
		log.Println("b.db.GetStartmessageId startMenu error: ", err)
	}
	b.cleareMessage(userId, startMessageId)

	// Установка нового startMessage.
	if err := b.db.UpdateStartMessageId(userId, newStartMessage.MessageID); err != nil {
		return err
	}

	return nil
}

// Отправка adMessage.
func (b *BotTelegram) sendAdMessage(userId int64) error {
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
	newAdMessage, err := b.bot.Send(botMsg)
	if err != nil {
		return fmt.Errorf("error send new adMessage: %w", err)
	}

	// Сохранение adMessageId.
	if err := b.db.AddUserMessageId(userId, newAdMessage.MessageID); err != nil {
		return err
	}

	// Удаление если возможно старого startMessage.
	adMessageId, err := b.db.GetAdMessageId(userId)
	if err != nil {
		log.Println("b.db.GetStartmessageId startMenu error: ", err)
	}
	b.cleareMessage(userId, adMessageId)

	// Установка нового adMessage.
	if err := b.db.UpdateAdMessageId(userId, newAdMessage.MessageID); err != nil {
		return err
	}

	return nil
}

// Очистка кэшей пользователя.
// TODO удалить существующий кэш пользователя.
func (b *BotTelegram) clearCacheOfUser(userId int64) error {
	delete(b.adEventCreatingCache, userId)
	delete(b.adEventsCache, userId)
	return nil
}
