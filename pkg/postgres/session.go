package postgres

import (
	"ada-telegram-bot/pkg/models"
	"encoding/json"
	"fmt"
)

// Сохранение сессии
func (t *TelegramBotDB) SaveSassion(userId int64, s *models.Session) (uuid string, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Получение json сессии
	jsS, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("error encoding session: %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO public.%s (user_id, session) values ($1, $2) RETURNING uuid;`, sessionsTable)
	if err := tx.QueryRow(query, userId, jsS).Scan(&uuid); err != nil {
		return "", fmt.Errorf("error insert session: %w", err)
	}

	return uuid, nil
}

// Получение крайней добавленной сессии
func (t *TelegramBotDB) GetLastSession(userId int64) (s *models.Session, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var session models.Session
	query := fmt.Sprintf(`SELECT session FROM public.%s WHERE user_id=$1 ORDER BY create_at DESC LIMIT 1;`, sessionsTable)
	if err := tx.QueryRow(query, userId).Scan(&session); err != nil {
		return nil, fmt.Errorf("error select last session: %w", err)
	}

	return &session, nil
}

// Удаление крайней добавленной сессии
func (t *TelegramBotDB) DeleteLastSession(userId int64) error {
	var err error
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`DELETE FROM public.%s WHERE uuid=(SELECT uuid FROM public.%s WHERE user_id=$1 ORDER BY create_at DESC LIMIT 1);`, sessionsTable)
	if _, err := tx.Exec(query, userId); err != nil {
		return fmt.Errorf("error delete last session: %w", err)
	}

	return nil
}

// Удаление сессий пользователя
func (t *TelegramBotDB) DeleteAllSession(userId int64) error {
	var err error
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`DELETE FROM public.%s WHERE uuid=(SELECT uuid FROM public.%s WHERE user_id=$1);`, sessionsTable)
	if _, err := tx.Exec(query, userId); err != nil {
		return fmt.Errorf("error delete last session: %w", err)
	}

	return nil
}
