package telegram

import (
	"ada-telegram-bot/pkg/models"
	"ada-telegram-bot/pkg/service"
	"ada-telegram-bot/pkg/subscriber"
	"fmt"
	"time"
)

func createStaticsBriefText(d *models.DataForStatistics) string {
	return fmt.Sprintf(`
	<b>📈 Статистика</b>
<b>Продано реклам:</b> %d
<b>Куплено реклам:</b> %d
<b>Кол-во взаимных пиаров:</b> %d
<b>Кол-во бартеров:</b> %d
<b>Прибыль:</b> %d
<b>Траты:</b> %d
<b>Чистая прибыль:</b> %d
`, d.CountAdEventSale, d.CountAdEventBuy, d.CountAdEventMutaul, d.CountAdEventBarter, d.Profit, d.Losses, d.Profit-d.Losses)
}

// Создание текст-описания ad события.
func createTextAdEventDescription(a *models.AdEvent) (descriptionAdEvent string) {
	// TODO Требуется избавиться от данного решения к 2065 году.
	if a.DateEnd == "02.01.06 15:04" {
		a.DateEnd = ""
	}

	if a.SubscribersOfChannel == 0 {
		a.SubscribersOfChannel, _ = subscriber.Parse(a.Channel)
	}

	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Тип:</b> <u>продажа рекламы</u>
		- <b>Покупатель:</b> %s
		- <b>Канал покупателя:</b> %s
		- <b>Кол-во подписчиков на канале покупателя:</b> %d
		- <b>Стоимость:</b> %d
		- <b>Дата размещения:</b> %s
		- <b>Дата удаления:</b> %s`, a.Partner, a.Channel, a.SubscribersOfChannel, a.Price, a.DateStart, a.DateEnd)
		if a.ArrivalOfSubscribers != 0 {
			descriptionAdEvent = descriptionAdEvent + fmt.Sprintf(`
			-<b>Приход подписчиков:</b> %d`, a.ArrivalOfSubscribers)
		}
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Тип:</b> <u>покупка рекламы</u>
		- <b>Продавец:</b> %s
		- <b>Канал продавца:</b> %s
		- <b>Кол-во подписчиков на канале продавца:</b> %d
		- <b>Стоимость:</b> %d
		- <b>Дата размещения:</b> %s`, a.Partner, a.Channel, a.SubscribersOfChannel, a.Price, a.DateStart)
		if a.ArrivalOfSubscribers != 0 {
			descriptionAdEvent = descriptionAdEvent + fmt.Sprintf(`
			-<b>Приход подписчиков:</b> %d`, a.ArrivalOfSubscribers)
		}
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Тип:</b> <u>взаимный пиар</u>
		- <b>Партнер:</b> %s
		- <b>Канал партнера:</b> %s
		- <b>Кол-во подписчиков на канале партнера:</b> %d
		- <b>Стоимость:</b> %d
		- <b>Дата размещения:</b> %s
		- <b>Дата удаления:</b> %s`, a.Partner, a.Channel, a.SubscribersOfChannel, a.Price, a.DateStart, a.DateEnd)
		if a.ArrivalOfSubscribers != 0 {
			descriptionAdEvent = descriptionAdEvent + fmt.Sprintf(`
			-<b>Приход подписчиков:</b> %d`, a.ArrivalOfSubscribers)
		}
	case models.TypeBarter:
		descriptionAdEvent = fmt.Sprintf(`
		- <b>Тип:</b> <u>бартер</u>
		- <b>Партнер:</b> %s
		- <b>Канал партнера:</b> %s
		- <b>Кол-во подписчиков на канале партнера:</b> %d
		- <b>Стоимость:</b> %d
		- <b>Дата размещения:</b> %s
		- <b>Дата удаления:</b> %s`, a.Partner, a.Channel, a.SubscribersOfChannel, a.Price, a.DateStart, a.DateEnd)
		if a.ArrivalOfSubscribers != 0 {
			descriptionAdEvent = descriptionAdEvent + fmt.Sprintf(`
			-<b>Приход подписчиков:</b> %d`, a.ArrivalOfSubscribers)
		}
	}

	descriptionAdEvent = descriptionAdEvent + "\n"

	return descriptionAdEvent
}

// Создание текста оповещения для размещения рекламы.
func createTextAlertForAdEventPosting(a *models.AdEvent, minutesLeftAlert int64) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s Вы должны разместить рекламу. Подробнее:
		`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s Ваша реклама будет размещена. Подробнее:
		`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s у Вас начнется взаимный пиар. Подробнее:
		`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeBarter:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s Вы должны разместить бартер. Подробнее:
		`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	}

	return descriptionAdEvent
}

// Создание текста оповещения для удаления рекламы.
func createTextAlertForAdEventDelete(a *models.AdEvent, minutesLeftAlert int64) (descriptionAdEvent string) {
	switch a.Type {
	case models.TypeSale:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s Вы можете удалить рекламу. Подробнее:`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeBuy:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s Ваша реклама будет удалена. Подробнее:`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeMutual:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s у Вас закончится взаимный пиар. Подробнее:`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	case models.TypeBarter:
		descriptionAdEvent = fmt.Sprintf(`
		❗️ Через %s у Вас закончится бартер. Подробнее:`+createTextAdEventDescription(a), getTextTime(minutesLeftAlert))
	}

	return descriptionAdEvent
}

// Получение правильного текста в зависимости от времени.
func getTextTime(minutes int64) string {
	var textTime string
	if minutes/60 < 1 {
		// Минуты
		if minutes == 1 {
			textTime = fmt.Sprintf("<b>%d</b> минута", minutes)
		} else if minutes >= 2 && minutes <= 4 {
			textTime = fmt.Sprintf("<b>%d</b> минуты", minutes)
		} else {
			textTime = fmt.Sprintf("<b>%d</b> минут", minutes)
		}
	} else {
		// Часы
		hours := minutes / 60
		switch {
		case hours == 1 || hours == 21:
			textTime = fmt.Sprintf("<b>%d</b> час", hours)
		case hours >= 2 && hours <= 4 || hours >= 22 && hours <= 24:
			textTime = fmt.Sprintf("<b>%d</b> часа", hours)
		default:
			textTime = fmt.Sprintf("<b>%d</b> часов", hours)
		}
	}

	return textTime
}

// Возвращает пример даты.
func getTextExampleDate() (string, error) {
	date, err := service.ParseTimeToUserDate(time.Now())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`
	В данный момент бот использует только время по МСК 'UTC+3'.
	<b>Пример:</b> <code>%s</code> `, date), nil
}

// Пример ссылки канала.
func getExampleChannelUrl() string {
	return `<b>Пример:</b> <code>@DmitriySergeevich22</code> или <code>https://t.me/DmitriySergeevich22</code>`
}

// Пример ссылки партнера.
func getExamplePartnerUrl() string {
	return `<b>Пример:</b> <code>@DmitriiSergeevich22</code> или <code>https://t.me/DmitriiSergeevich22</code>`
}

// Текст получение стоимости события.
func textForGetPrice(t models.TypeAdEvent) (string, error) {
	switch t {
	case models.TypeSale:
		return "✍️ Теперь требуется отправить стоимость рекламы.\n<b>Пример:</b> <code>1000</code>", nil
	case models.TypeBuy:
		return "✍️ Теперь требуется отправить стоимость рекламы.\n<b>Пример:</b> <code>1000</code>", nil
	case models.TypeMutual:
		return `✍️ Теперь требуется отправить стоимость поста взаимного пиара.
<b>Пример:</b> <code>1000</code>
Можно указать <code>-1000</code> если была доплата с Вашей стороны.`, nil
	case models.TypeBarter:
		return `✍️ Теперь требуется отправить прибыль с бартера.
<b>Пример:</b> <code>1000</code>. <code>0</code> - eсли считать прибыль не требуется.`, nil
	default:
		return "", fmt.Errorf("unknow type adEvent")
	}
}

// Текст обновления стоимости события.
func textForUpdatePrice() string {
	return "✍️ Требуется отправить новую стоимость.\n<b>Пример:</b> <code>1000</code>"
}

// Текст получения url канала.
func textForGetDateChannelUrl(t models.TypeAdEvent) (string, error) {
	switch t {
	case models.TypeSale:
		return "✍️ Теперь требуется отправить ссылку на рекламируемый Вами канал.\n" + getExampleChannelUrl(), nil
	case models.TypeBuy:
		return "✍️ Теперь требуется отправить ссылку на канал, в котором выйдет Ваша реклама.\n" + getExampleChannelUrl(), nil
	case models.TypeMutual:
		return "✍️ Теперь требуется отправить ссылку на канал, с которым будет взаимный пиар.\n" + getExampleChannelUrl(), nil
	case models.TypeBarter:
		return "✍️ Теперь требуется отправить ссылку на канал/магазин партнера по бартеру.\n" + getExampleChannelUrl(), nil
	default:
		return "", fmt.Errorf("unknow type adEvent. typeEvent: %s", t)
	}
}

// Текст получения даты размещения.
func textForGetDateStart(t models.TypeAdEvent) (string, error) {
	exampleDate, err := getTextExampleDate()
	if err != nil {
		return "", err
	}

	switch t {
	case models.TypeSale:
		return "✍️ Теперь требуется отправить дату и время размещения рекламы." + exampleDate, nil
	case models.TypeBuy:
		return "✍️ Теперь требуется отправить дату и время размещения рекламы." + exampleDate, nil
	case models.TypeMutual:
		return "✍️ Теперь требуется отправить дату и время размещения поста." + exampleDate, nil
	case models.TypeBarter:
		return "✍️ Теперь требуется отправить дату и время размещения поста." + exampleDate, nil
	default:
		return "", fmt.Errorf("unknow type adEvent. typeEvent: %s", t)
	}
}

// Текст получения даты удаления.
func textForGetDateEnd(t models.TypeAdEvent) (string, error) {
	exampleDate, err := getTextExampleDate()
	if err != nil {
		return "", err
	}

	switch t {
	case models.TypeSale:
		return "✍️ Теперь требуется отправить дату и время удаления рекламы." + exampleDate, nil
	case models.TypeBuy:
		return "✍️ Теперь требуется отправить дату и время удаления рекламы." + exampleDate, nil
	case models.TypeMutual:
		return "✍️ Теперь требуется отправить дату и время удаления поста." + exampleDate, nil
	case models.TypeBarter:
		return "✍️ Теперь требуется отправить дату и время удаления поста." + exampleDate, nil
	default:
		return "", fmt.Errorf("unknow type adEvent. typeEvent: %s", t)
	}
}

// Тест успешного добавления даты и времени удаления.
func textForSuccessfullyAddDeleteDate(t models.TypeAdEvent) (string, error) {
	switch t {
	case models.TypeSale:
		return "🎉 <b>Дата и время удаления рекламы добавлены!</b>", nil
	case models.TypeBuy:
		return "🎉 <b>Дата и время удаления рекламы добавлены!</b>", nil
	case models.TypeMutual:
		return "🎉 <b>Дата и время удаления поста добавлены!</b>", nil
	case models.TypeBarter:
		return "🎉 <b>Дата и время удаления поста добавлены!</b>", nil
	default:
		return "", fmt.Errorf("unknow type adEvent. typeEvent: %s", t)
	}
}
