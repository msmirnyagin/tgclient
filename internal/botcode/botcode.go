package botcode

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendLog(token string, chatID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	bot.Debug = false
	if chatID == 0 {
		return nil
	}
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
	return nil
}

func GetCode(token string) (string, int64, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return "", 0, err
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "accepted")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			return update.Message.Text, update.Message.Chat.ID, nil
		}
	}
	return "", 0, nil
}
