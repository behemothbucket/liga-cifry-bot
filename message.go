package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
	"strings"
)

func handleMessage(update tgbotapi.Update) {
	if update.Message.From == nil {
		return
	}

	user := update.Message.From
	text := update.Message.Text

	logMessage(user, text)

	var err error

	if isValidMessageText(text) {
		if isAuthorized {
			if strings.HasPrefix(text, "/") {
				err = handleCommand(update, text)
			} else if searchMode {
				if len(searchCriterias) == 0 {
					err = sendFinalMessage(update.Message.Chat.ID)
					searchMode = false
				} else {
					err = sendAcceptMessage(update.Message.Chat.ID)
					delete(searchCriterias, currentCriteria)
				}
			} else {
				err = sendMenuMessage(update.Message.Chat.ID)
			}
		} else if text == "В клюве" {
			isAuthorized = true
			enabledInlineKeyboard = true
			enabledKeyboard = false
			err = sendMenuMessage(update.Message.Chat.ID)
		} else {
			err = sendAuthorizationMessage(update.Message.Chat.ID)
		}
	} else {
		err = sendMediaErrorMessage(update.Message.Chat.ID)
	}

	if err != nil {
		log.Println("Ошибка:", err.Error())
	}
}

func isValidMessageText(message string) bool {
	if reflect.TypeOf(message).Kind() == reflect.String && message != "" {
		return true
	} else {
		return false
	}
}

func logMessage(user *tgbotapi.User, text string) {
	userName := user.UserName
	firstName := user.FirstName
	lastName := user.LastName
	userID := user.ID

	if lastName != "" {
		lastName = " " + lastName
	}

	log.Printf("https://t.me/%s [%d] (%s%s) написал(а) '%s'", userName, userID, firstName, lastName, text)
}

func sendAcceptMessage(chatID int64) error {
	return sendFormattedMessage(chatID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
}

func sendFinalMessage(chatID int64) error {
	err := sendFormattedMessage(chatID, "✅ Вы заполнили все критерии")
	if err != nil {
		return err
	}
	err = sendMenuMessage(chatID)
	return err
}

func sendAuthorizationMessage(chatID int64) error {
	return sendFormattedMessage(chatID, "❗️Пожалуйста, войдите в группу по ссылке-приглашению")
}

func sendMediaErrorMessage(chatID int64) error {
	return sendFormattedMessage(chatID, "❌ Файлы, фото/видео и другие медиа <b>не принимаются</b>")
}

func sendFormattedMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if enabledInlineKeyboard && !searchMode {
		msg.ReplyMarkup = getMainMenuMarkup()
	}

	if enabledKeyboard {
		msg.ReplyMarkup = authorizationKeyboardButton
	}

	_, err := bot.Send(msg)

	return err
}

func sendMenuMessage(chatID int64) error {
	err := SendMenu(chatID)
	return err
}
