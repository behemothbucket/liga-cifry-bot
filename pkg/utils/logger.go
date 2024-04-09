package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LogMessage(message *tgbotapi.Message) {
	userName := message.From.UserName
	firstName := message.From.FirstName
	lastName := message.From.LastName
	userID := message.From.ID
	text := message.Text
	chatID := message.Chat.ID
	var groupName string

	if lastName != "" {
		lastName = " " + lastName
	}

	if message.Chat.Title != "" {
		groupName = message.Chat.Title
	}

	log.Printf("https://t.me/%s [ID:%d] (%s%s) send message '%s' to chat [chatID:%d, group:%s]",
		userName, userID, firstName, lastName, text, chatID, groupName)
}
