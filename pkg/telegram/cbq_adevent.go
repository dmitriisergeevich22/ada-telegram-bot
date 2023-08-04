package telegram

import (
	"ada-telegram-bot/pkg/models"
	"ada-telegram-bot/pkg/service"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CBQ AdEventCreate

func cbqAdEventCreate(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "<b>✍️ Выберите тип события:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продажа рекламы", "ad_event.create.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покупка рекламы", "ad_event.create.buy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар", "ad_event.create.mutual"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Бартер", "ad_event.create.barter"),
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

func cbqAdEventCreateSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Type:      models.TypeSale,
	}
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := "✍️ Теперь требуется отправить ссылку на чат с покупателем.\n" + getExamplePartnerUrl()
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateSale: %w", err)
	}

	return nil
}

func cbqAdEventCreateBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Type:      models.TypeBuy,
	}
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := "✍️ Теперь требуется отправить ссылку на продавца.\n" + getExamplePartnerUrl()
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateBuy: %w", err)
	}

	return nil
}

func cbqAdEventCreateMutual(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Type:      models.TypeMutual,
	}
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := "✍️ Теперь требуется отправить ссылку на пратнера по взаимному пиару.\n" + getExamplePartnerUrl()
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateMutual: %w", err)
	}

	return nil
}

func cbqAdEventCreateBarter(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Создание кэша ad события.
	adEvent := models.AdEvent{
		UserId:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05.999"),
		Type:      models.TypeBarter,
	}
	b.adEventCreatingCache[userId] = &adEvent

	b.db.SetStepUser(userId, "ad_event.create.partner")

	text := "✍️ Теперь требуется отправить ссылку на пратнера по бартеру.\n" + getExamplePartnerUrl()
	botMsg := tgbotapi.NewEditMessageText(userId, messageId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateBarter: %w", err)
	}

	return nil
}

func cbqAdEventCreateEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	// Валидация события.
	if !fullDataAdEvent(adEvent) {
		botMsg := tgbotapi.NewMessage(userId, "Были введены не все данные, что бы повторить воспользуйтесь командой <b>/start</b>")
		botMsg.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню.", "start"),
			),
		)
		botMsg.ReplyMarkup = keyboard

		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Сохранение события в бд.
	_, err = b.db.AdEventCreation(adEvent)
	if err != nil {
		return err
	}

	// Отправка сообщения.
	text := "<b>🎊 Отлично! Событие добавлено!</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventCreateEnd: %w", err)
	}

	// Очистка кэша.
	delete(b.adEventCreatingCache, userId)
	return nil
}

// CBQ AdEventView

func cbqAdEventView(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	text := "<b>✍️ Выберите тип событий:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Проданная реклама", "ad_event.view.sale"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Купленная реклама", "ad_event.view.buy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Взаимный пиар", "ad_event.view.mutual"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Бартер", "ad_event.view.barter"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("Все типы", "ad_event.view.any"),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventView: %w", err)
	}

	return nil
}

func cbqAdEventViewAny(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())+";any;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeLastMonth())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "ad_event.view.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)+";any;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAny: %w", err)
	}

	return nil
}

func cbqAdEventViewSale(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())+";sale;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())+";sale;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())+";sale;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeLastMonth())+";sale;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())+";sale;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())+";sale;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "ad_event.view.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)+";sale;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewSale: %w", err)
	}

	return nil
}

func cbqAdEventViewBuy(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())+";buy;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())+";buy;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())+";buy;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeLastMonth())+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())+";buy;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())+";buy;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "ad_event.view.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)+";buy;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAny: %w", err)
	}

	return nil
}

func cbqAdEventViewMutual(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())+";mutual;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())+";mutual;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeLastMonth())+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "ad_event.view.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)+";mutual;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAny: %w", err)
	}

	return nil
}

func cbqAdEventViewBarter(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := "<b>🕐 Выберите период:</b>"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вчера", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeYesterday())+";barter;1"),
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeToday())+";barter;1"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeTomorrow())+";barter;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeLastMonth())+";barter;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisMonth())+";barter;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий год", "ad_event.view.select?"+service.ParseTimesToRangeDate(service.GetTimeRangeThisYear())+";barter;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За все время", "ad_event.view.select?"+service.ParseTimesToRangeDate(models.MinTime, models.MaxTime)+";barter;1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewBarter: %w", err)
	}

	return nil
}

func cbqAdEventViewSelect(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID
	lenRow := 3

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// TODO Сохранение в кэш выборки ( удалить как появится сохранение в БД)
	if !b.toCache(userId, "cbqAdEventViewSelectData", cbqData) {
		return fmt.Errorf("cbqAdEventViewSelect: error save cbqData in cache")
	}

	// Парсинг данных.
	data, err := parseDataAdEventView(cbqData)
	if err != nil {
		return err
	}

	// Получение событий из БД.
	adEvents, err := b.db.GetRangeAdEventsOfUser(userId, data.TypeAdEvent, data.StartDate, data.EndDate)
	if err != nil {
		return err
	}

	// Разбиение событий и сохранение в кэш.
	b.adEventsCache[userId] = service.ChunkSlice(adEvents, lenRow)

	// Отображение событий.
	text, keyboard, err := createTextAndKeyboardForAdEventView(b, userId, data)
	if err != nil {
		return err
	}

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.DisableWebPagePreview = true
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAnyAll: %w", err)
	}

	return nil
}

func parseDataAdEventView(cbqData string) (data *models.CbqDataForCbqAdEventViewSelect, err error) {
	// ad_event.view.any.select?14.05.2023 00:00;14.05.2023 23:59;any;1
	dataSlice := strings.Split(cbqData, ";")
	if len(dataSlice) != 4 {
		return nil, fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}
	data = new(models.CbqDataForCbqAdEventViewSelect)

	data.StartDate, err = service.ParseUserDateToTime(dataSlice[0])
	if err != nil {
		return nil, err
	}

	data.EndDate, err = service.ParseUserDateToTime(dataSlice[1])
	if err != nil {
		return nil, err
	}

	data.TypeAdEvent = models.TypeAdEvent(dataSlice[2])

	pageForDisplay, err := strconv.Atoi(dataSlice[3])
	if err != nil {
		return nil, fmt.Errorf("error pasge PageForDisplay: %w", err)
	}
	data.PageForDisplay = pageForDisplay

	return data, nil
}

func createTextAndKeyboardForAdEventView(b *BotTelegram, userId int64, data *models.CbqDataForCbqAdEventViewSelect) (string, tgbotapi.InlineKeyboardMarkup, error) {
	lenRow := 3

	adEvents, err := b.getAdEventsCache(userId)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, err
	}

	if len(adEvents) == 0 {
		text := `<b>🗓 Нет событий.</b>`
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view."+string(data.TypeAdEvent)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
			),
		)

		return text, keyboard, nil
	}

	// Создание кнопок.
	text := fmt.Sprintf(`<b>🗓 Выбранные события. Страница %d/%d. </b>
	✔️ Выберите номер события на <b>кнопках ниже</b> для редактирования события.
	`, data.PageForDisplay, len(adEvents))

	bufButtonRows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	bufButtonRow := make([]tgbotapi.InlineKeyboardButton, 0, lenRow)
	for i, adEvent := range adEvents[data.PageForDisplay-1] {
		buttonId := fmt.Sprintf("%d", i+1)
		buttonData := fmt.Sprintf("ad_event.control?%d", adEvent.Id)
		button := tgbotapi.NewInlineKeyboardButtonData(buttonId, buttonData)
		bufButtonRow = append(bufButtonRow, button)

		text = text + fmt.Sprintf("\n<b>    ✍️ Событие № %s</b>:", buttonId)
		text = text + createTextAdEventDescription(&adEvent)
	}
	bufButtonRows = append(bufButtonRows, bufButtonRow)

	if len(adEvents) > 1 {
		pageRow := createPageRowForViewAdEvent(data, len(adEvents))
		bufButtonRows = append(bufButtonRows, pageRow)
	}

	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view."+string(data.TypeAdEvent)),
	)
	bufButtonRows = append(bufButtonRows, backRow)

	startMenuRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
	)
	bufButtonRows = append(bufButtonRows, startMenuRow)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(bufButtonRows...)

	return text, keyboard, nil
}

func createPageRowForViewAdEvent(data *models.CbqDataForCbqAdEventViewSelect, maxPage int) []tgbotapi.InlineKeyboardButton {
	buffButton := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	if data.PageForDisplay-1 > 0 {
		textDataPreviousPage := fmt.Sprintf("ad_event.view.select?%s;%s;%d",
			service.ParseTimesToRangeDate(data.StartDate, data.EndDate), data.TypeAdEvent, data.PageForDisplay-1)
		buffButton = append(buffButton, tgbotapi.NewInlineKeyboardButtonData("<<", textDataPreviousPage))
	}

	if data.PageForDisplay+1 <= maxPage {
		textDataNextPage := fmt.Sprintf("ad_event.view.select?%s;%s;%d",
			service.ParseTimesToRangeDate(data.StartDate, data.EndDate), data.TypeAdEvent, data.PageForDisplay+1)
		buffButton = append(buffButton, tgbotapi.NewInlineKeyboardButtonData(">>", textDataNextPage))
	}

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view.any"),
	)

	return tgbotapi.NewInlineKeyboardRow(buffButton...)
}

// CBQ AdEventControl

func cbqAdEventControl(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	text := "📝 Выберите действие:"

	deleteButtonData := fmt.Sprintf("ad_event.delete?%d", adEventId)
	updatePartnerButtonData := fmt.Sprintf("ad_event.update.partner?%d", adEventId)
	updateChannelButtonData := fmt.Sprintf("ad_event.update.channel?%d", adEventId)
	updatePriceButtonData := fmt.Sprintf("ad_event.update.price?%d", adEventId)
	dateStartButtonData := fmt.Sprintf("ad_event.update.date_start?%d", adEventId)
	dateEndButtonData := fmt.Sprintf("ad_event.update.date_end?%d", adEventId)
	arrivalOfSubscribersButtonData := fmt.Sprintf("ad_event.update.arrival_of_subscribers?%d", adEventId)

	cbqAdEventViewSelectData, ok := b.fromCache(userId, "cbqAdEventViewSelectData")
	if !ok {
		return fmt.Errorf("cbqAdEventControl: error getCache")
	}
	cbqAdEventViewSelectDataString, ok := cbqAdEventViewSelectData.(string)
	if !ok {
		return fmt.Errorf("cbqAdEventControl: error parse cbqAdEventViewSelectData to string")
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить", deleteButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить ссылку на партнера", updatePartnerButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить ссылку на канал партнера", updateChannelButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить стоимость", updatePriceButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить дату и время размещения рекламы", dateStartButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить дату и время удаления рекламы", dateEndButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Внести приход подписчиков", arrivalOfSubscribersButtonData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "ad_event.view.select?"+cbqAdEventViewSelectDataString),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	if err := b.sendMessage(userId, tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventView: %w", err)
	}

	return nil
}

// CBQ AdEventDelete

func cbqAdEventDelete(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	aE, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}

	text := "<b>⚠️ Вы точно хотите удалить событие?</b>" + createTextAdEventDescription(aE)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "ad_event.delete.end?"+strconv.Itoa(int(adEventId))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEventId)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)

	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAny: %w", err)
	}
	return nil
}

func cbqAdEventDeleteEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Удаление события.
	if err := b.db.AdEventDelete(adEventId); err != nil {
		return err
	}

	text := "❌ Событие удалено! ❌"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.DisableWebPagePreview = true
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventViewAnyAll: %w", err)
	}

	return nil
}

// CBQ AdEventUpdate

func cbqAdEventUpdatePartner(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.partner"); err != nil {
		return err
	}

	text := "✍️ Требуется отправить новую ссылку на партнера.\n" + getExamplePartnerUrl()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}

	return nil
}

func cbqAdEventUpdateChannel(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.channel"); err != nil {
		return err
	}

	text := "✍️ Требуется отправить новую ссылку на канал.\n" + getExampleChannelUrl()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}

	return nil
}

func cbqAdEventUpdatePrice(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.price"); err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, textForUpdatePrice(), keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}

	return nil
}

func cbqAdEventUpdateDateStart(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.date_start"); err != nil {
		return err
	}

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	text := `✍️ Теперь требуется отправить дату и время размещения рекламы.` + exampleDate
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}
	return nil
}

func cbqAdEventUpdateDateEnd(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.date_end"); err != nil {
		return err
	}

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	text := `✍️ Теперь требуется отправить новую дату и время удаления рекламы.` + exampleDate
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}

	return nil
}

func cbqAdEventUpdateArrivalOfSubscribers(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Получение данных cbq.
	_, cbqData, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	// Парсинг данных cbq.
	adEventId, err := cbqParseDataGetAdEventId(cbqData)
	if err != nil {
		return err
	}

	// Добавление события в кэш.
	adEvent, err := b.db.GetAdEvent(adEventId)
	if err != nil {
		return err
	}
	b.adEventCreatingCache[userId] = adEvent

	// Установка шага.
	if err := b.db.SetStepUser(userId, "ad_event.update.arrival_of_subscribers"); err != nil {
		return err
	}

	text := `✍️ Требуется отправить приход подписчиков:
	<b>Пример:</b> 1000`
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqAdEventUpdatePartner: %w", err)
	}

	return nil
}
