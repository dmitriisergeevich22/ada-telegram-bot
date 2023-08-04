package subscriber

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func Parse(url string) (subscription int64, err error) {
	switch ParseTypeUrl(url) {
	case TelegramUrlType:
		return getCurrentSubscriptionFromTelegramChannel(url)
	}

	log.Println("Unkown URL type: ", url)
	return 0, nil
}

func ParseTypeUrl(url string) UrlType {
	if url[:13] == "https://t.me/" {
		return TelegramUrlType
	}
	return UnknowUrl
}

// Получение кол-ва подписчиков телеграм канала.
func getCurrentSubscriptionFromTelegramChannel(url string) (subscription int64, err error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error send request in url (%s): %w", url, err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading body response from url (%s): %w", url, err)
	}

	subscribers, err := scanSubscribersInTelegram(doc)
	if err != nil {
		return 0, err
	}

	return subscribers, nil
}

// Сканирования кол-ва подписчиков в телеграме.
func scanSubscribersInTelegram(doc *goquery.Document) (int64, error) {
	nameClass := doc.Find(".tgme_page_extra")
	if nameClass == nil {
		return 0, fmt.Errorf("scanSubscribersInTelegram: nil nameClass")
	}

	dataHtml, err := nameClass.Html()
	if err != nil {
		return 0, fmt.Errorf("scanSubscribersInTelegram: error getting html: %w", err)
	}

	var subscribersStrings []rune
	for _, r := range dataHtml {
		if r >= '0' && r <= '9' {
			subscribersStrings = append(subscribersStrings, r)
		}
	}

	subscribers, _ := strconv.ParseInt(string(subscribersStrings), 0, 64)

	return subscribers, nil
}
