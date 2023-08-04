package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func cbqInfo(b *BotTelegram, cbq *tgbotapi.CallbackQuery) error {
	userId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	// Сборка сообщения.
	text := `
ℹ️ <b>Информация о проекте</b> ℹ️

Данный бот был создан для помощи в работе с рекламными интеграциями. 
	
<b>Что умеет бот на текущий момент:</b> 
• <u>Создавать</u> рекламные интеграции:
  - Покупка рекламы
  - Продажа рекламы
  - Взаимный пиар
  - Бартер
• <u>Напоминать</u> о предстоящих рекламных интеграциях
• <u>Отслеживать</u> финансовые показатели
`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "start"),
		),
	)
	botMsg := tgbotapi.NewEditMessageTextAndMarkup(userId, messageId, text, keyboard)
	botMsg.ParseMode = tgbotapi.ModeHTML

	if err := b.sendMessage(userId, botMsg); err != nil {
		return fmt.Errorf("error edit msg in cbqHelpInfo: %w", err)
	}

	return nil
}
