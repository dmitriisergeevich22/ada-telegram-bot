package postgres

import (
	"ada-telegram-bot/pkg/models"
	"database/sql"
	"fmt"
)


// Создание default пользователя.
// Если пользователь уже создан - ошибка равна nil.
// UserId == ChatId
func (t *TelegramBotDB) CreateUser(userId int64, userUrl, firstName string) (err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Default User
	dU := models.User{
		Id:           userId,
		Name:         firstName,
		UserURL:      "@" + userUrl,
		Step:         "start",
		Login:        userUrl,
		PasswordHash: "123",
	}

	// Создание default пользователя.
	sql := fmt.Sprintf(`INSERT INTO public.%s (id, name, user_url, step, login, password)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING;`, usersTable)
	if _, err := tx.Exec(sql, dU.Id, dU.Name, dU.UserURL, dU.Step, dU.Login, dU.PasswordHash); err != nil {
		return fmt.Errorf("error create default user. err: %w", err)
	}

	return nil
}

// Получение данных пользователя.
func (t *TelegramBotDB) GetUser(userId int64) (user *models.User, err error) {
	tx := t.db.MustBegin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	u := new(models.User)
	query := fmt.Sprintf(`SELECT (id, created_at, name, user_url, step, login, password)
	FROM public.%s WHERE id=$1;`, usersTable)
	if err := tx.QueryRow(query, userId).Scan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("error scan user data. err: %w", err)
	}

	return u, nil
}
