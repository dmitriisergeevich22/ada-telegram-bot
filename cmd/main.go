package main

import (
	"AdaTelegramBot/internal/postgresql"
	"AdaTelegramBot/internal/telegram"
	"log"
	"os"

	"github.com/spf13/viper"
)

func main() {
	// Отображение переменных.
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	log.Println("SSL_MODE:", os.Getenv("SSL_MODE"))
	log.Println("TG_TOKEN:", os.Getenv("TG_TOKEN"))

	// Инициализация конфигурации
	if err := initConfig(); err != nil {
		log.Println("main: error initConfig: ", err)
		return
	}

	// Подключение к БД.
	db := postgresql.NewDB()

	// Инициализация телеграмм бота.
	tgBot, err := telegram.NewBotTelegram(postgresql.NewTelegramBotDB(db))
	if err != nil {
		log.Println("main: error telegram.NewBotTelegram: ", err)
		return
	}

	// Запуск бота.
	if err := tgBot.StartBotUpdater(); err != nil {
		log.Println("main: error tgBot.StartBotUpdater: ", err)
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
