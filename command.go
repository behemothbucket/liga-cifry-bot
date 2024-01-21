package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func handleCommand(update tgbotapi.Update, command string) error {
	var err error

	switch command {
	case "/start":
		if isAuthorized {
			err = SendMenu(update.Message.Chat.ID)
		} else {
			err = sendAuthorizationMessage(update.Message.Chat.ID)
		}
	case "/login":
		err = sendAuthorizationMessage(update.Message.Chat.ID)
	}

	return err
}
