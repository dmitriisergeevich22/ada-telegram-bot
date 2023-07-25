package ada

import (
	"ada-telegram-bot/pkg/models"
	"ada-telegram-bot/pkg/service"
	"ada-telegram-bot/pkg/subscriber"
	"fmt"
	"log"
	"math"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// Оповещение о предстоящих событиях.
func (a *AdaBot) alertTicker() error {
	timeAlert := viper.GetInt("ada_bot.speed_chek_alert")
	if timeAlert == 0 {
		timeAlert = 10
	}

	for {
		time.Sleep(time.Duration(timeAlert) * time.Second)
		if err := handlerAlertsTick(a); err != nil {
			log.Println(err)
		}
	}
}

// Отправка оповещения пользователю.
func sendAlertMessage(a *AdaBot, userId int64, c tgbotapi.Chattable) error {
	botMsg, err := a.bot.Send(c)
	if err != nil {
		return err
	}

	// Добавления ID сообщения в бд.
	if err := a.db.AddUserMessageId(userId, botMsg.MessageID, "alert"); err != nil {
		return err
	}

	// Обновление даты оповещения.
	if err := a.db.UpdateTimeLastAlert(userId, time.Now()); err != nil {
		return err
	}

	return nil
}

func handlerAlertsTick(a *AdaBot) error {
	timeStart, _ := service.GetTimeRangeToday()
	_, timeEnd := service.GetTimeRangeTomorrow()
	cashAdEvents, err := a.db.GetRangeAdEvents(models.TypeAny, timeStart, timeEnd)
	if err != nil {
		return fmt.Errorf("handlerAlertsTick: error GetRangeAdEvents: %w", err)
	}

	for _, aE := range cashAdEvents {
		// Проврека последнего оповещения.
		timeLastAlert, err := a.db.GetTimeLastAlert(aE.UserId)
		if err != nil {
			return fmt.Errorf("handlerAlertsTick: error GetTimeLastAlert: %w", err)
		}

		// Оповещение не чаще чем раз в 1 минуту.
		// TODO: Добавить метку об успешном оповещении события. Так как может быть 2 события разный в 1 минуту.
		if int64(math.Abs(time.Since(timeLastAlert).Minutes())) > 1 {
			if err := aletrPosting(a, &aE); err != nil {
				return fmt.Errorf("handlerAlertsTick: error aletrPosting: %w", err)
			}
			if err := aletrDelete(a, &aE); err != nil {
				return fmt.Errorf("handlerAlertsTick: error aletrDelete: %w", err)
			}
		}
	}

	return nil
}

// Оповещение о размещении рекламы.
func aletrPosting(a *AdaBot, aE *models.AdEvent) error {
	timeDateStart, err := service.ParseUserDateToTime(aE.DateStart)
	if err != nil {
		return fmt.Errorf("aletrPosting: error ParseUserDateToTime: %w", err)
	}

	// Событие прошло.
	if time.Since(timeDateStart).Minutes() > 0 {
		return nil
	}

	// Сохранение подписчиков на момент выхода рекламы.
	if int64(math.Abs(time.Since(timeDateStart).Minutes())) == 0 {
		currentSub, err := subscriber.Parse(aE.Channel)
		if err != nil {
			return fmt.Errorf("aletrPosting: error subscriber_parser.Parse: %w", err)
		}

		if err := a.db.UpdatePartnerChannelSubscribersInStart(aE.Id, currentSub); err != nil {
			return fmt.Errorf("aletrPosting: error UpdatePartnerChannelSubscribersInStart: %w", err)
		}

		// TODO сохранить кол-во подписчиков канала пользователя.

		if err := a.db.UpdateTimeLastAlert(aE.UserId, time.Now()); err != nil {
			return fmt.Errorf("aletrPosting: error UpdateTimeLastAlert: %w", err)
		}
	}

	timeLeft := int64(math.Abs(time.Since(timeDateStart).Minutes()))
	if checkTimeAlert(aE.UserId, timeLeft) {
		text := createTextAlertForAdEventPosting(aE, timeLeft)
		botMsg := tgbotapi.NewMessage(aE.UserId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		botMsg.DisableWebPagePreview = true
		if err := sendAlertMessage(a, aE.UserId, botMsg); err != nil {
			return fmt.Errorf("aletrPosting: error sendAlertMessage: %w", err)
		}
		log.Println("aletrPosting: successfully send posting alert: ", aE)
	}

	return nil
}

// Оповещение о удалении рекламы.
func aletrDelete(a *AdaBot, aE *models.AdEvent) error {
	timeDateEnd, err := service.ParseUserDateToTime(aE.DateEnd)
	if err != nil {
		return fmt.Errorf("aletrDelete: error ParseUserDateToTime: %w", err)
	}

	// Событие прошло.
	if time.Since(timeDateEnd).Minutes() > 0 {
		return nil
	}

	// Сохранение подписчиков на момент завершения рекламы.
	if int64(math.Abs(time.Since(timeDateEnd).Minutes())) == 0 {
		currentSub, err := subscriber.Parse(aE.Channel)
		if err != nil {
			return fmt.Errorf("aletrDelete: error subscriber_parser.Parse: %w", err)
		}

		if err := a.db.UpdatePartnerChannelSubscribersInEnd(aE.Id, currentSub); err != nil {
			return fmt.Errorf("aletrDelete: error UpdatePartnerChannelSubscribersInEnd: %w", err)
		}

		// TODO сохранить кол-во подписчиков канала пользователя.

		if err := a.db.UpdateTimeLastAlert(aE.UserId, time.Now()); err != nil {
			return fmt.Errorf("aletrDelete: error UpdateTimeLastAlert: %w", err)
		}
	}

	timeLeft := int64(math.Abs(time.Since(timeDateEnd).Minutes()))
	// Удаления  отображаются только за 1 час.
	if timeLeft > 60 {
		return nil
	}

	if checkTimeAlert(aE.UserId, timeLeft) && aE.Type != models.TypeBuy && aE.Type != models.TypeBarter {
		text := createTextAlertForAdEventDelete(aE, timeLeft)
		botMsg := tgbotapi.NewMessage(aE.UserId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		botMsg.DisableWebPagePreview = true
		if err := sendAlertMessage(a, aE.UserId, botMsg); err != nil {
			return fmt.Errorf("aletrDelete: error sendAlertMessage: %w", err)
		}

		log.Println("aletrDelete: successfully send delete alert: ", aE)
	}

	return nil
}

// Проверка доступа к оповещениям
func checkTimeAlert(userId, timeLeft int64) bool {
	// var timeAlerts []int64
	// TODO Смотрим на какое время установил предупреждения полульзователь.
	timeAlerts := []int64{1440, 60, 30, 10}

	for _, timeAlert := range timeAlerts {
		if timeAlert == timeLeft {
			return true
		}
	}

	return false
}