package botcode

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetCode(token string) (string, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return "", err
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			return update.Message.Text, nil
		}
	}
	return "", nil
}
