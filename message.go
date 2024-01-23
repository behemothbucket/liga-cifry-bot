package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
)

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

func (b *Bot) sendAcceptMessage(chatID int64) error {
	return b.sendMarkupMessage(chatID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
}

func (b *Bot) sendSearchFinalMessage(chatID int64) error {
	err := b.sendMessage(chatID, "✅ Вы заполнили все критерии")
	if err != nil {
		return err
	}
	showSearchResultsMode = false
	err = b.sendMainMenu(chatID)
	return err
}

func (b *Bot) sendAuthorizationMessage(chatID int64) error {
	return b.sendMarkupMessage(chatID, "❗️Пожалуйста, войдите в группу по ссылке-приглашению")
}

func (b *Bot) sendMediaErrorMessage(chatID int64) error {
	return b.sendMessage(chatID, "❌ Файлы, фото/видео и другие медиа <b>не принимаются</b>")
}

// TODO передавать какой-то объект для определения нужности markup
func (b *Bot) sendMarkupMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if searchMode {
		msg.ReplyMarkup = cancelMenuMarkup
	}

	if enabledInlineKeyboard && !searchMode && !showSearchResultsMode {
		msg.ReplyMarkup = mainMenuMarkup
	}

	if enabledKeyboard {
		msg.ReplyMarkup = authorizationKeyboardButton
	}

	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := b.bot.Send(msg)
	return err
}
