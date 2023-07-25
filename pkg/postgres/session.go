package postgres

import "ada-telegram-bot/pkg/models"

// Добавление сессии
func (t *TelegramBotDB) SaveSassion(userId int64, s *models.Session) (uuid string, err error) {
	return "uuid", nil
}

// Получение крайней добавленной сессии
func (t *TelegramBotDB) GetLastSession(userId int64) (s *models.Session, err error) {
	return nil, nil
}

// Удаление крайней добавленной сессии
func (t *TelegramBotDB) DeleteLastSession(userId int64) error {
	return nil
}

// Удаление сессий пользователя
func (t *TelegramBotDB) DeleteAllSession(userId int64) error {
	return nil
}