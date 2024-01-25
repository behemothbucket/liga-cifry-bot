package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
)

func (b *Bot) isValidMessageText(update tgbotapi.Update) bool {
	message := update.Message
	if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
		return true
	} else {
		return false
	}
}

func logMessage(user *tgbotapi.User, text string, chatID int64) {
	userName := user.UserName
	firstName := user.FirstName
	lastName := user.LastName
	userID := user.ID

	if lastName != "" {
		lastName = " " + lastName
	}

	log.Printf("https://t.me/%s [ID:%d] (%s%s) написал(а) '%s' в чат [chatID:%d]", userName, userID, firstName, lastName, text, chatID)
}

func (b *Bot) sendAcceptMessage(chatID int64) error {
	return b.sendMarkupMessage(chatID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
}

func (b *Bot) sendSearchFinalMessage(chatID int64) error {
	err := b.SendMessage(chatID, "✅ Вы заполнили все критерии")
	if err != nil {
		return err
	}
	showSearchResultsMode = false
	err = b.sendMainMenu(chatID)
	return err
}

func (b *Bot) sendMediaErrorMessage(chatID int64) error {
	return b.SendMessage(chatID, "❌ Файлы, фото/видео и другие медиа <b>не принимаются</b>")
}

// TODO передавать какой-то объект для определения нужности markup
func (b *Bot) sendMarkupMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if searchMode {
		msg.ReplyMarkup = cancelMenuMarkup
	}

	if !searchMode && !showSearchResultsMode {
		msg.ReplyMarkup = mainMenuMarkup
	}

	if showSearchResultsMode {
		msg.ReplyMarkup = backToMainMenuMarkup
		showSearchResultsMode = false
	}

	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := b.bot.Send(msg)
	return err
}
