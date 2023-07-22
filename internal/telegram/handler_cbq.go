package telegram

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func parseCbq(cbq *tgbotapi.CallbackQuery) (path []string, data string, err error) {
	cbqDataSlice := strings.Split(cbq.Data, "?")
	if len(cbqDataSlice) < 1 {
		return nil, "", fmt.Errorf("len cbq incorrect. cbq: %s ", cbq.Data)
	}

	cbqPathSlice := strings.Split(cbqDataSlice[0], ".")
	if len(cbqPathSlice) < 1 {
		return nil, "", fmt.Errorf("len cbq path incorrect. cbq: %s ", cbq.Data)
	}

	switch len(cbqDataSlice) {
	case 1:
		return cbqPathSlice, "", nil
	case 2:
		return cbqPathSlice, cbqDataSlice[1], nil
	default:
		return nil, "", fmt.Errorf("len cbq incorrect. cbq: %s", cbq.Data)
	}
}

func (b *BotTelegram) handlerCbq(cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	fmt.Printf("Info %s; user=%s; CBQ=%s;\n", time.Now().Format("2006-01-02 15:04:05.999"), cbq.From.UserName, cbq.Data)
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch path[0] {
	case "start":
		if err := b.cmdStart(cbq.Message); err != nil {
			log.Println("error in cmdStart: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "ad_event":
		if err := handlerCbqAdEvent(b, cbq); err != nil {
			log.Println("error in handlerCbqAdEvent: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "statistics":
		if err := handlerCbqStatistics(b, cbq); err != nil {
			log.Println("error in handlerCbqStatistics: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "info":
		if err := handlerCbqInfo(b, cbq); err != nil {
			log.Println("error in handlerCbqInfo: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	case "help":
		if err := handlerCbqHelp(b, cbq); err != nil {
			log.Println("error in handlerCbqHelp: ", err)
			b.sendRequestRestartMsg(userId)
			return err
		}
	}

	return nil
}

func handlerCbqAdEvent(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
	case "ad_event.create":
		if err := cbqAdEventCreate(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.sale":
		if err := cbqAdEventCreateSale(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.buy":
		if err := cbqAdEventCreateBuy(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.mutual":
		if err := cbqAdEventCreateMutual(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.barter":
		if err := cbqAdEventCreateBarter(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.end":
		if err := cbqAdEventCreateEnd(b, cbq); err != nil {
			return err
		}
	case "ad_event.create.price.skip":
		b.cleareMessage(userId, messageId)

		b.db.SetStepUser(userId, "ad_event.create.date_start")
		adEvent, err := b.getAdEventCreatingCache(userId)
		if err != nil {
			return err
		}

		text, err := textForGetDateStart(adEvent.Type)
		if err != nil {
			return err
		}
		botMsg := tgbotapi.NewMessage(userId, text)
		botMsg.ParseMode = tgbotapi.ModeHTML
		if err := b.sendMessage(userId, botMsg); err != nil {
			return err
		}
	case "ad_event.create.date_end.skip":
		b.cleareMessage(userId, messageId)

		if err := adEventCreateLastMessage(b, userId); err != nil {
			return err
		}
	case "ad_event.view":
		if err := cbqAdEventView(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.any":
		if err := cbqAdEventViewAny(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.sale":
		if err := cbqAdEventViewSale(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.buy":
		if err := cbqAdEventViewBuy(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.mutual":
		if err := cbqAdEventViewMutual(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.barter":
		if err := cbqAdEventViewBarter(b, cbq); err != nil {
			return err
		}
	case "ad_event.view.select":
		if err := cbqAdEventViewSelect(b, cbq); err != nil {
			return err
		}
	case "ad_event.control":
		if err := cbqAdEventControl(b, cbq); err != nil {
			return err
		}
	case "ad_event.delete":
		if err := cbqAdEventDelete(b, cbq); err != nil {
			return err
		}
	case "ad_event.delete.end":
		if err := cbqAdEventDeleteEnd(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.partner":
		if err := cbqAdEventUpdatePartner(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.channel":
		if err := cbqAdEventUpdateChannel(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.price":
		if err := cbqAdEventUpdatePrice(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.date_start":
		if err := cbqAdEventUpdateDateStart(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.date_end":
		if err := cbqAdEventUpdateDateEnd(b, cbq); err != nil {
			return err
		}
	case "ad_event.update.arrival_of_subscribers":
		if err := cbqAdEventUpdateArrivalOfSubscribers(b, cbq); err != nil {
			return err
		}
	}
	return nil
}

func handlerCbqStatistics(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
	case "statistics":
		if err := cbqStatistics(b, cbq); err != nil {
			return err
		}
	case "statistics.brief":
		if err := cbqStatisticsBrief(b, cbq); err != nil {
			return err
		}
	case "statistics.brief.select":
		if err := cbqStatisticsBriefSelect(b, cbq); err != nil {
			return err
		}
	}
	return nil
}

func handlerCbqInfo(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
	case "info":
		if err := cbqInfo(b, cbq); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func handlerCbqHelp(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	path, _, err := parseCbq(cbq)
	if err != nil {
		return err
	}

	switch strings.Join(path, ".") {
	case "help":
		if err := cbqHelp(b, cbq); err != nil {
			return err
		}
	case "help.feature":
		if err := cbqHelpFeature(b, cbq); err != nil {
			return err
		}
	case "help.error":
		if err := cbqHelpError(b, cbq); err != nil {
			return err
		}
	}

	return nil
}
