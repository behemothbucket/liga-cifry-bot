package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
	"strings"
)

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	if user == nil {
		return
	}

	userName := user.UserName
	firstName := user.FirstName
	lastName := user.LastName
	userID := user.ID

	if lastName != "" {
		lastName = " " + lastName
	}

	log.Printf("https://t.me/%s [%d] (%s%s) написал(а) '%s'", userName, userID, firstName, lastName, text)

	var err error

	if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
		if strings.HasPrefix(text, "/") {
			err = handleCommand(message.Chat.ID, text)
		} else if searchMode {
			if len(searchCriterias) == 0 {
				msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Вы заполнили все критерии</b>✅")
				msg.ParseMode = tgbotapi.ModeHTML
				_, err = bot.Send(msg)
				err = SendMenu(message.Chat.ID)
				searchMode = false
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
			msg.ParseMode = tgbotapi.ModeHTML
			_, err = bot.Send(msg)
			delete(searchCriterias, currentCriteria)
		} else {
			err = SendMenu(message.Chat.ID)
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Файлы, фото/видео и другие медиа не принимаются ❌")
		_, err = bot.Send(msg)
		err = SendMenu(message.Chat.ID)
	}

	if err != nil {
		log.Printf("Ошибка: %s", err.Error())
	}
}
