package telegram

import (
	"ada-telegram-bot/pkg/models"
	"ada-telegram-bot/pkg/subscriber"
	"ada-telegram-bot/pkg/service"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик сообщений.
func (b *BotTelegram) handlerMessage(msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	fmt.Printf("Info %s: user=%s; MSG=%s;\n", time.Now().Format("2006-01-02 15:04:05.999"), msg.From.UserName, msg.Text)
	step, err := b.db.GetStepUser(userId)
	if err != nil {
		return err
	}

	// Сообщение обрабатываеются отталкиваясь от текущего шага пользователя.
	switch step {
	case "ad_event.create.partner":
		if err := adEventPartner(b, msg); err != nil {
			log.Println("error in adEventPartner: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.channel":
		if err := adEventChannel(b, msg); err != nil {
			log.Println("error in adEventChannel: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.price":
		if err := adEventPrice(b, msg); err != nil {
			log.Println("error in adEventPrice: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.date_start":
		if err := adEventDateStart(b, msg); err != nil {
			log.Println("error in adEventDateStart: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.create.date_end":
		if err := adEventDateEnd(b, msg); err != nil {
			log.Println("error in adEventDateEnd: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.partner":
		if err := adEventUpdatePartner(b, msg); err != nil {
			log.Println("error in adEventUpdatePartner: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.channel":
		if err := adEventUpdateChannel(b, msg); err != nil {
			log.Println("error in adEventUpdateChannel: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.price":
		if err := adEventUpdatePrice(b, msg); err != nil {
			log.Println("error in adEventUpdatePrice: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.date_start":
		if err := adEventUpdateDateStart(b, msg); err != nil {
			log.Println("error in adEventUpdateDateStart: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.date_end":
		if err := adEventUpdateDateEnd(b, msg); err != nil {
			log.Println("error in adEventUpdateDateEnd: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event.update.arrival_of_subscribers":
		if err := adEventUpdateArrivalOfSubscribers(b, msg); err != nil {
			log.Println("error in adEventUpdateArrivalOfSubscribers: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	default:
		botMsg := tgbotapi.NewMessage(userId, "Не получается обработать сообщение... 😔")
		botMsg.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
			),
		)
		botMsg.ReplyMarkup = keyboard

		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}

	}

	return nil
}

func adEventPartner(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы отправили некорректную ссылку на пользователя, попробуйте снова.\n"+getExamplePartnerUrl())
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if models.RegxUrlType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	adEvent.Partner = msg.Text
	b.db.SetStepUser(userId, "ad_event.create.channel")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на пользователя добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	// Получение канала.
	text, err := textForGetDateChannelUrl(adEvent.Type)
	if err != nil {
		return err
	}
	botMsg = tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventChannel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID
	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы отправили некорректную ссылку на канал, попробуйте снова."+getExampleChannelUrl())
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Приведение в единный тип.
	if models.RegxUrlType2.MatchString(msg.Text) {
		msg.Text = "https://t.me/" + msg.Text[1:]
	}

	subChannel, err := subscriber.Parse(msg.Text)
	if err != nil {
		return err
	}
	adEvent.SubscribersOfChannel = subChannel

	adEvent.Channel = msg.Text
	b.db.SetStepUser(userId, "ad_event.create.price")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на канал добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	// Получение стоимости.
	text, err := textForGetPrice(adEvent.Type)
	if err != nil {
		return err
	}
	botMsg = tgbotapi.NewMessage(userId, text)
	if adEvent.Type == models.TypeMutual || adEvent.Type == models.TypeBarter {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пропустить", "ad_event.create.price.skip"),
			),
		)
		botMsg.ReplyMarkup = keyboard
	}
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventPrice(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную стоимость, попробуйте снова.
		<b>Пример:</b> <code>1000</code>`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	price, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}

	adEvent.Price = price
	b.db.SetStepUser(userId, "ad_event.create.date_start")

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Стоимость добавлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	text, err := textForGetDateStart(adEvent.Type)
	if err != nil {
		return err
	}
	botMsg = tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventDateStart(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateStart = msg.Text

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время размещения рекламы добавлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	// Отправка сообщения об получении даты удаления.
	text, err := textForGetDateEnd(adEvent.Type)
	if err != nil {
		return err
	}
	b.db.SetStepUser(userId, "ad_event.create.date_end")
	botMsg = tgbotapi.NewMessage(userId, text)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "ad_event.create.date_end.skip"),
		),
	)
	botMsg.ReplyMarkup = keyboard
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventDateEnd(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Заполнение информации в хэш-таблице ad событий.
	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateEnd = msg.Text

	// Сравнение даты размещения и удаления.
	durationDateStart, err := service.ParseUserDateToTime(adEvent.DateStart)
	if err != nil {
		return fmt.Errorf("error parse durationDateStart: %w", err)
	}

	durationDateEnd, err := service.ParseUserDateToTime(adEvent.DateEnd)
	if err != nil {
		return fmt.Errorf("error parse durationDateEnd: %w", err)
	}

	if durationDateEnd.Sub(durationDateStart) <= 0 {
		botMsg := tgbotapi.NewMessage(userId, "Вы ввели дату удаления поста меньше даты размещения поста, попробуйте снова.")
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	// Сообщение об успешном добавлении.
	text, err := textForSuccessfullyAddDeleteDate(adEvent.Type)
	if err != nil {
		return err
	}
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	// Отправка завершающего создания ad события сообщения.
	if err := adEventCreateLastMessage(b, userId); err != nil {
		return err
	}

	return nil
}

func adEventCreateLastMessage(b *BotTelegram, userId int64) error {
	aE, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	text := "<b>✍️ Вы хотите создать данное событие?</b>"
	text = text + createTextAdEventDescription(aE)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "ad_event.create.end"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", "start"),
		),
	)
	botMsg := tgbotapi.NewMessage(userId, text)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.ReplyMarkup = keyboard
	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}
	return nil
}

func adEventUpdatePartner(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы отправили некорректную ссылку на партнера, попробуйте снова."+getExamplePartnerUrl())
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.Partner = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на партнера обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateChannel(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxUrlType1.MatchString(msg.Text) && !models.RegxUrlType2.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, "Вы отправили некорректную ссылку на канал, попробуйте снова."+getExampleChannelUrl())
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.Channel = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Ссылка на канал обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdatePrice(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxPrice.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную стоимость, попробуйте снова.
		<b>Пример:</b> <code>1000</code>`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	price, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}

	adEvent.Price = price

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Стоимость обновлена!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateDateStart(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateStart = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время размещения рекламы обновлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateDateEnd(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	exampleDate, err := getTextExampleDate()
	if err != nil {
		return err
	}
	if !models.RegxAdEventDate.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректную дату и время, попробуйте снова.
		`+exampleDate)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}
	adEvent.DateEnd = msg.Text

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Дата и время удаления рекламы обновлены!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}

func adEventUpdateArrivalOfSubscribers(b *BotTelegram, msg *tgbotapi.Message) error {
	userId := msg.Chat.ID

	if !models.RegxArrivalOfSubscribers.MatchString(msg.Text) {
		botMsg := tgbotapi.NewMessage(userId, `Вы отправили некорректный приход подписчиков, попробуйте снова.
		<b>Пример:</b> 1000`)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
		return nil
	}

	adEvent, err := b.getAdEventCreatingCache(userId)
	if err != nil {
		return err
	}

	arrivalOfSubscribers, err := strconv.ParseInt(msg.Text, 0, 64)
	if err != nil {
		return err
	}
	adEvent.ArrivalOfSubscribers = arrivalOfSubscribers

	if err := b.db.AdEventUpdate(adEvent); err != nil {
		return err
	}

	botMsg := tgbotapi.NewMessage(userId, "🎉 <b>Приход подписчиков обновлен!</b>")
	botMsg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("ad_event.control?%d", adEvent.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg.ReplyMarkup = keyboard

	if err := b.sendMessage(userId, botMsg); err != nil {
		return err
	}

	return nil
}
