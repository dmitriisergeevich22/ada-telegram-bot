package models

import "time"

type DB interface {
	Ping() error
	Close() error
	UserDB
	SessionDB
	AdEventDB
	MessageDB
	StatisticsDB
}

type UserDB interface {
	// Создание пользователя.
	CreateUser(chatId int64, userUrl, firstName string) error
	// Получение данных пользователя.
	GetUser(userId int64) (user *User, err error)
	// Получение последней даты оповещения.
	GetTimeLastAlert(userId int64) (timeLastAlert time.Time, err error)
	// Обновление последней даты оповещения.
	UpdateTimeLastAlert(userId int64, timeLastAlert time.Time) error
	// Обновление даты последней активности пользователя.
	UpdateLastActive(userId int64) error
	// 
}

type SessionDB interface {
	// Сохранение состояния сессии
	SaveSassion(userId int64, s *Session) (uuid string, err error)
	// Получение крайней сессии
	GetLastSession(userId int64) (s *Session, err error)
	// Удаление крайней сессии
	DeleteLastSession(userId int64) error
	// Удалить все сессии
	DeleteAllSession(userId int64) error
}

type AdEventDB interface {
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
	// Обновление кол-ва подписчиков канала на момент выхода рекламы.
	UpdatePartnerChannelSubscribersInStart(adEventId, subscribers int64) error
	// Обновление кол-ва подписчиков канала на момент завершения рекламы.
	UpdatePartnerChannelSubscribersInEnd(adEventId, subscribers int64) error
}

type MessageDB interface {
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
}

type StatisticsDB interface {
	// Получение данных пользователя для статистики.
	GetRangeDataForStatistics(userId int64, typeAdEvent TypeAdEvent, startDate, endDate time.Time) (data *DataForStatistics, err error)
}
