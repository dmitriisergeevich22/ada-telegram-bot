package models

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Пользователь при регистрации.
type User struct {
	Id             int64  `json:"id"`                                   // Chat_ID
	CreatedAt      string `json:"createdAt" db:"created_at"`            // Дата создания.
	Name           string `json:"name" db:"name"`                       // Имя пользователя.
	UserURL        string `json:"userUrl" db:"user_url"`                // Ссылка пользователя.
	Step           string `json:"stap" db:"stap"`                       // Шаг пользвателя (на каком шаге находится пользователь)
	Login          string `json:"login" db:"login"`                     // Логин пользователя.
	PasswordHash   string `json:"password" db:"password"`               // Хэш пользовательского пароля.
	StartMessageId int    `json:"startMessageId" db:"start_message_id"` // Индефикатор сообщения startMessage.
	AdMessageId    int    `json:"adMessageId" db:"ad_message_id"`       // Индефикатор сообщения adMessage.
	InfoMessageId  int    `json:"infoMessageId" db:"info_message_id"`   // Индефикатор сообщения infoMessage.
	DateLastAlert  string `json:"DateLastAlert" db:"date_last_alert"`   // Даты последнего оповещения пользователя.
}

var (
	MinTime = time.Date(1969, time.January, 1, 0, 0, 0, 0, time.UTC)      // Минимальная дата time.Time
	MaxTime = time.Date(2068, time.December, 31, 23, 59, 59, 0, time.UTC) // Максимальная дата time.Time
	MinDate = "01.01.69 00:00"                                            // 1969-01-01 00:00:00 +0300 MSK
	MaxDate = "31.12.68 23:59"                                            // 2068-12-31 23:59:00 +0300 MSK

	ErrUserNotFound = fmt.Errorf("user not found")
	// Example: "22.08.22 16:30"
	RegxAdEventDate = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01]).(0[1-9]|1[0-2]).(\d{2}) ([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	// Example: "https://t.me/nikname", "https://www.instagram.com/nikname.store/"
	RegxUrlType1 = regexp.MustCompile(`^https://[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+(/[a-zA-Z0-9-]*)*`)
	// Example: "@nikname"
	RegxUrlType2 = regexp.MustCompile(`^@[a-zA-Z0-9_]+$`)
	// Example: 1000
	RegxPrice = regexp.MustCompile(`[0-9]+`)
	// Example: 1000
	RegxArrivalOfSubscribers = regexp.MustCompile(`[0-9]+`)
	// Example: 1
	RegxId = regexp.MustCompile(`[0-9]+`)
)

// Типы CallbackQuery.

type CbqStatic tgbotapi.CallbackQuery  // CallbackQuery без CbqData
type CbqDinamic tgbotapi.CallbackQuery // CallbackQuery с CbqData
type CbqPath []string                  // Путь
type CbqData []byte

// Тип рекламных интегаций.
type TypeAdEvent string

const (
	TypeAny    TypeAdEvent = "any"
	TypeSale   TypeAdEvent = "sale"
	TypeBuy    TypeAdEvent = "buy"
	TypeMutual TypeAdEvent = "mutual"
	TypeBarter TypeAdEvent = "barter"
)

type AdEvent struct {
	// Индификатор события.
	Id int64 `json:"id" db:"id"`
	// Дата создания события.
	CreatedAt string `json:"createdAt" db:"created_at"`
	// Id пользователя.
	UserId int64 `json:"userId" db:"user_id"`
	// Тип события.
	Type TypeAdEvent `json:"type" db:"type"`
	// Ссылка партнера.
	Partner string `json:"partner" db:"partner"`
	// Ссылка на канал партнера.
	Channel string `json:"channel" db:"channel"`
	// Подписчики канала.
	SubscribersOfChannel int64 `json:"subscribersOfChannel" db:"subscribers_of_channel"`
	// Стоимость.
	Price int64 `json:"price" db:"price"`
	// Дата начала события. "02.01.06 15:04"
	DateStart string `json:"dateStart" db:"date_start"`
	// Дата завершения события. "02.01.06 15:04"
	DateEnd string `json:"dateEnd" db:"date_end"`
	// Приход подписчиков.
	ArrivalOfSubscribers int64 `json:"arrivalOfSubscribers" db:"arrival_of_subscribers"`
	// Кол-во подписчиков в начале события.
	SubscribersInStart int64 `json:"subscribersInStart" db:"subscribers_in_start"`
	// Кол-во подписчиков в конце события.
	SubscribersInEnd int64 `json:"subscribersInEnd" db:"subscribers_in_end"`
	// Кол-во подписчиков партнера в начале события.
	PartnerChannelSubscribersInStart int64 `json:"partnerChannelSubscribersInStart" db:"partner_channel_subscribers_in_start"`
	// Кол-во подписчиков партнера в конце события.
	PartnerChannelSubscribersInEnd int64 `json:"partnerChannelSubscribersInEnd" db:"partner_channel_subscribers_in_end"`
}

// Данные для создания статистики.
type DataForStatistics struct {
	CountAdEventSale   int64 // Кол-во проданных реклам.
	CountAdEventBuy    int64 // Кол-во купленных реклам.
	CountAdEventMutaul int64 // Кол-во взаимных пиаров.
	CountAdEventBarter int64 // Кол-во бартеров.
	Profit             int64 // Прибыль.
	Losses             int64 // Убытки.
}

// Сессия пользователя.
type OldSession struct {
	DomainPath string                 // Наименование основной цепочки.
	Step       int64                  // Шаг в цепочке.
	StateMsg   string                 // Состояние ожидающих данных в Msg.
	Cache      map[string]interface{} // Кэш сессии.
}

// БД для телеграмм бота.
type TelegramBotDB interface {
	Close() error

	// Получение данных пользователя.
	GetUserData(userId int64) (user *User, err error)
	// Создание пользователя.
	DefaultUserCreation(chatId int64, userUrl, firstName string) error
	// Получение последней даты оповещения.
	GetTimeLastAlert(userId int64) (timeLastAlert time.Time, err error)
	// Обновление последней даты оповещения.
	UpdateTimeLastAlert(userId int64, timeLastAlert time.Time) error
	// Обновление кол-ва подписчиков канала на момент выхода рекламы.
	UpdatePartnerChannelSubscribersInStart(adEventId, subscribers int64) error
	// Обновление кол-ва подписчиков канала на момент завершения рекламы.
	UpdatePartnerChannelSubscribersInEnd(adEventId, subscribers int64) error

	// Получение ad события.
	GetAdEvent(adEventId int64) (*AdEvent, error)
	// Получение всех ad событий в указаном диапазоне времени.
	GetRangeAdEvents(typeAdEvent TypeAdEvent, startDate, endDate time.Time) ([]AdEvent, error)
	// Получение всех ad событий пользователя запрашиваемого типа.
	GetAdEventsOfUser(userId int64, typeAdEvent TypeAdEvent) ([]AdEvent, error)
	// Получение всех ad событий пользователя запрашиваемого типа в указаном диапазоне времени.
	GetRangeAdEventsOfUser(userId int64, typeAdEvent TypeAdEvent, startDate, endDate time.Time) ([]AdEvent, error)
	// Создание ad события.
	AdEventCreation(adEvent *AdEvent) (int64, error)
	// Удаление ad события.
	AdEventDelete(eventId int64) error
	// Обновление информации о приходе подписчиков.
	AdEventUpdate(adEvent *AdEvent) error
	// Установка шага пользователя.
	SetStepUser(userId int64, step string) error
	// Получение текущего шага пользователя.
	GetStepUser(userId int64) (step string, err error)

	// Добавление messageId пользователя.
	AddUserMessageId(userId int64, messageId int) error
	// Удаление messageId пользователя.
	DeleteUsermessageId(messageId int) error
	// Возвращает список messageIds пользователя.
	GetUserMessageIds(userId int64) ([]int, error)
	// Возвращает startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	GetStartMessageId(userId int64) (messageId int, err error)
	// Обновление startmessageId. Это сообщение которое не удаляется а меняется на меню команды /start.
	UpdateStartMessageId(userId int64, messageId int) error
	// Возвращает admessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	GetAdMessageId(userId int64) (messageId int, err error)
	// Обновление AdmessageId. Это сообщение которое не удаляется, купленная в боте реклама.
	UpdateAdMessageId(userId int64, messageId int) error
	// Обновление даты последней активности пользователя.
	UpdateLastActive(userId int64) error

	// Получение данных пользователя для статистики.
	GetRangeDataForStatistics(userId int64, typeAdEvent TypeAdEvent, startDate, endDate time.Time) (data *DataForStatistics, err error)
}

type CbqDataForCbqAdEventViewSelect struct {
	StartDate      time.Time   // Начальная дата событий.
	EndDate        time.Time   // Конечная дата событий.
	TypeAdEvent    TypeAdEvent // Тип событий.
	PageForDisplay int         // Страница для отображения.
}
