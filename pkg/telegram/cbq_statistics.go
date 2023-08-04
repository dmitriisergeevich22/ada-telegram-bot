package telegram

import (
	"ada-telegram-bot/pkg/models"
	"ada-telegram-bot/pkg/service"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "📈 <b>Статистика:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Краткая статистика", "statistics.brief"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Полная статистика", "statistics.full"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqStatistics: %w", err)
	}

	return nil
}

func cbqStatisticsBrief(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущая неделя", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisWeek())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "statistics.brief.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "statistics.brief.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "statistics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreate: %w", err)
	}

	return nil
}

func cbqStatisticsBriefSelect(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	_, data, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	dataSlice := strings.Split(data, ";")
	if len(dataSlice) != 2 {
		return fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}

	startDate, err := service.ParseUserDateToTime(dataSlice[0])
	if err != nil {
		return err
	}
	endDate, err := service.ParseUserDateToTime(dataSlice[1])
	if err != nil {
		return err
	}

	// Получение данных из БД.
	d, err := b.db.GetRangeDataForStatistics(userId, models.TypeAny, startDate, endDate)
	if err != nil {
		return err
	}

	// Создание краткой статистики.
	text := createStaticsBriefText(d)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "statistics.brief"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqStatisticsBriefSelect: %w", err)
	}

	return nil
}
