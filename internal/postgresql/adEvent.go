package postgresql

import (
	"AdaTelegramBot/internal/models"
	"time"

	"database/sql"
	"fmt"
)

func (t *TelegramBotDB) GetAdEvent(adEventId int64) (adEvent *models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
	arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
	FROM public.%s WHERE id=$1;`, adEventsTable)

	var aE models.AdEvent
	if err := tx.QueryRow(query, adEventId).Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
		&aE.DateStart, &aE.DateEnd, &aE.ArrivalOfSubscribers, &aE.PartnerChannelSubscribersInStart,
		&aE.PartnerChannelSubscribersInEnd); err != nil {
		return nil, fmt.Errorf("GetAdEvent: error scan AdEvent in GetAdEvent: %w", err)
	}

	aE.DateStart, err = dbDateToUserDate(aE.DateStart)
	if err != nil {
		return nil, fmt.Errorf("GetAdEvent: error dbDateToUserDate: %w", err)
	}

	aE.DateEnd, err = dbDateToUserDate(aE.DateEnd)
	if err != nil {
		return nil, fmt.Errorf("GetAdEvent: error dbDateToUserDate: %w", err)
	}

	return &aE, nil
}

func (t *TelegramBotDB) GetRangeAdEvents(typeAdEvent models.TypeAdEvent, startTime, endTime time.Time) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	listAdEvent = make([]models.AdEvent, 0, 50)
	startDate, err := timeToDbDate(startTime)
	if err != nil {
		return nil, err
	}
	endDate, err := timeToDbDate(endTime)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE ((date_start BETWEEN $1 AND $2)
		OR (date_end BETWEEN $2 AND $1)) ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEvents: error select ad_events type `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE "type"=$3 AND ((date_start BETWEEN $1 AND $2)
		OR (date_end BETWEEN $1 AND $2)) ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, startDate, endDate, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEvents: error select ad_events type `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&aE.DateStart, &aE.DateEnd, &aE.ArrivalOfSubscribers,
			&aE.PartnerChannelSubscribersInStart, &aE.PartnerChannelSubscribersInEnd); err != nil {
			return nil, fmt.Errorf("GetRangeAdEvents: error scan AdEvent: %w", err)
		}

		aE.DateStart, err = dbDateToUserDate(aE.DateStart)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEvents: error dbDateToUserDate: %w", err)
		}

		aE.DateEnd, err = dbDateToUserDate(aE.DateEnd)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEvents: error dbDateToUserDate: %w", err)
		}

		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) GetAdEventsOfUser(userId int64, typeAdEvent models.TypeAdEvent) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	listAdEvent = make([]models.AdEvent, 0, 50)
	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE user_id=$1 ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId)
		if err != nil {
			return nil, fmt.Errorf("GetAdEventsOfUser: error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE user_id=$1 AND "type"=$2 ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("GetAdEventsOfUser: error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&aE.DateStart, &aE.DateEnd, &aE.ArrivalOfSubscribers,
			&aE.PartnerChannelSubscribersInStart, &aE.PartnerChannelSubscribersInEnd); err != nil {
			return nil, fmt.Errorf("GetAdEventsOfUser: error scan AdEvent: %w", err)
		}

		aE.DateStart, err = dbDateToUserDate(aE.DateStart)
		if err != nil {
			return nil, fmt.Errorf("GetAdEventsOfUser: error dbDateToUserDate: %w", err)
		}

		aE.DateEnd, err = dbDateToUserDate(aE.DateEnd)
		if err != nil {
			return nil, fmt.Errorf("GetAdEventsOfUser: error dbDateToUserDate: %w", err)
		}

		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) GetRangeAdEventsOfUser(userId int64, typeAdEvent models.TypeAdEvent, startTime, endTime time.Time) (listAdEvent []models.AdEvent, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	listAdEvent = make([]models.AdEvent, 0, 50)

	startDate, err := timeToDbDate(startTime)
	if err != nil {
		return nil, err
	}

	endDate, err := timeToDbDate(endTime)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if typeAdEvent == models.TypeAny {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE user_id=$1 AND ((date_start BETWEEN $2 AND $3)
		OR (date_end BETWEEN $2 AND $3)) ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEventsOfUser: error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	} else {
		query := fmt.Sprintf(`SELECT id, created_at, user_id, "type", partner, channel, price, date_start, date_end,
		arrival_of_subscribers, partner_channel_subscribers_in_start, partner_channel_subscribers_in_end
		FROM public.%s WHERE user_id=$1 AND "type"=$4 AND ((date_start BETWEEN $2 AND $3)
		OR (date_end BETWEEN $2 AND $3)) ORDER BY date_start ASC;`, adEventsTable)

		rows, err = tx.Query(query, userId, startDate, endDate, typeAdEvent)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEventsOfUser: error select ad_events TypeAny `%s`: %w", typeAdEvent, err)
		}
	}

	for rows.Next() {
		var aE models.AdEvent
		if err := rows.Scan(&aE.Id, &aE.CreatedAt, &aE.UserId, &aE.Type, &aE.Partner, &aE.Channel, &aE.Price,
			&aE.DateStart, &aE.DateEnd, &aE.ArrivalOfSubscribers,
			&aE.PartnerChannelSubscribersInStart, &aE.PartnerChannelSubscribersInEnd); err != nil {
			return nil, fmt.Errorf("GetRangeAdEventsOfUser: error scan AdEvent: %w", err)
		}

		aE.DateStart, err = dbDateToUserDate(aE.DateStart)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEventsOfUser: error aE.DateStart to dbDateToUserDate: %w", err)
		}

		aE.DateEnd, err = dbDateToUserDate(aE.DateEnd)
		if err != nil {
			return nil, fmt.Errorf("GetRangeAdEventsOfUser: error aE.DateEnd to dbDateToUserDate: %w", err)
		}

		listAdEvent = append(listAdEvent, aE)
	}

	return listAdEvent, nil
}

func (t *TelegramBotDB) AdEventCreation(aE *models.AdEvent) (eventId int64, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if aE.DateEnd == "" {
		aE.DateEnd = "02.01.06 15:04"
	}

	aE.DateStart, err = userDateToDbDate(aE.DateStart)
	if err != nil {
		return 0, fmt.Errorf("AdEventCreation: error parse aE.DateStart in userDateToDbDate: %w", err)
	}

	aE.DateEnd, err = userDateToDbDate(aE.DateEnd)
	if err != nil {
		return 0, fmt.Errorf("AdEventCreation: error parse aE.DateEnd in userDateToDbDate: %w", err)
	}

	sql := fmt.Sprintf(`INSERT INTO public.%s (user_id, "type", partner, channel, price, date_start, date_end)
	values ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, adEventsTable)
	if err := tx.QueryRow(sql, aE.UserId, aE.Type, aE.Partner, aE.Channel, aE.Price,
		aE.DateStart, aE.DateEnd).Scan(&eventId); err != nil {
		return 0, fmt.Errorf("error creation event: %w", err)
	}

	return eventId, nil
}

// TODO разделить на маленькие функции обновления.
func (t *TelegramBotDB) AdEventUpdate(aE *models.AdEvent) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if aE.Type != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET "type"=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Type, aE.Id); err != nil {
			return fmt.Errorf("error update type. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Partner != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET partner=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Partner, aE.Id); err != nil {
			return fmt.Errorf("error update partner. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Channel != "" {
		query := fmt.Sprintf(`UPDATE public.%s SET channel=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Channel, aE.Id); err != nil {
			return fmt.Errorf("error update channel. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.Price != 0 {
		query := fmt.Sprintf(`UPDATE public.%s SET price=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.Price, aE.Id); err != nil {
			return fmt.Errorf("error update price. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.DateStart != "" {
		aE.DateStart, err = userDateToDbDate(aE.DateStart)
		if err != nil {
			return fmt.Errorf("AdEventUpdate: error userDateToDbDate aE.DateStart: %w", err)
		}

		query := fmt.Sprintf(`UPDATE public.%s SET date_start=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.DateStart, aE.Id); err != nil {
			return fmt.Errorf("error update date_start. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.DateEnd != "" {
		aE.DateEnd, err = userDateToDbDate(aE.DateEnd)
		if err != nil {
			return fmt.Errorf("AdEventUpdate: error userDateToDbDate aE.DateEnd: %w", err)
		}

		query := fmt.Sprintf(`UPDATE public.%s SET date_end=$1
			WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.DateEnd, aE.Id); err != nil {
			return fmt.Errorf("error update date_end. eventId%d: %w", aE.Id, err)
		}
	}

	if aE.ArrivalOfSubscribers != 0 {
		query := fmt.Sprintf(`UPDATE public.%s SET arrival_of_subscribers=$1
		WHERE id=$2;`, adEventsTable)
		if _, err := tx.Exec(query, aE.ArrivalOfSubscribers, aE.Id); err != nil {
			return fmt.Errorf("error update arrival_of_subscribers. eventId%d: %w", aE.Id, err)
		}
	}

	return nil
}

// Удаление ad события.
func (t *TelegramBotDB) AdEventDelete(eventId int64) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	sql := fmt.Sprintf(`DELETE FROM public.%s WHERE id=$1;`, adEventsTable)
	if _, err := tx.Exec(sql, eventId); err != nil {
		return fmt.Errorf("error delete event: %w", err)
	}

	return nil
}

func (t *TelegramBotDB) UpdatePartnerChannelSubscribersInStart(adEventId, subscribers int64) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`UPDATE public.%s SET partner_channel_subscribers_in_start=$1
			WHERE id=$2;`, adEventsTable)
	if _, err := tx.Exec(query, subscribers, adEventId); err != nil {
		return fmt.Errorf("error update partner_channel_subscribers_in_start. eventId%d: %w", adEventId, err)
	}

	return nil
}

func (t *TelegramBotDB) UpdatePartnerChannelSubscribersInEnd(adEventId, subscribers int64) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`UPDATE public.%s SET partner_channel_subscribers_in_end=$1
			WHERE id=$2;`, adEventsTable)
	if _, err := tx.Exec(query, subscribers, adEventId); err != nil {
		return fmt.Errorf("error update partner_channel_subscribers_in_end. eventId%d: %w", adEventId, err)
	}

	return nil
}
