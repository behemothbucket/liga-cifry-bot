package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"reflect"
)

type Message struct {
	parseMode string
	chatID    int64
}

func isValidMessageText(update tgbotapi.Update) bool {
	var valid bool
	message := update.Message

	if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
		valid = true
	}

	return valid
}

func isJoinEvent(update tgbotapi.Update) bool {
	var event bool
	message := update.Message

	if len(message.NewChatMembers) != 0 {
		go func() {
			newUser := message.NewChatMembers[0]
			AddUserSql(newUser.ID, newUser.UserName, newUser.FirstName, newUser.LastName, newUser.IsBot)
		}()
		event = true
	}
	if update.Message.LeftChatMember != nil {
		go func() {
			userID := update.Message.LeftChatMember.ID
			DeleteUserSql(userID)
		}()
		event = true
	}

	return event
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

func (b *Bot) sendAcceptMessage(chatID int64) {
	b.SendMarkupMessage(chatID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
}

//func (b *Bot) sendSearchFinalMessage(chatID int64) {
//	b.SendMessage(chatID, "✅ Вы заполнили все критерии")
//	showSearchResultsMode = false
//	b.sendMainMenu(chatID)
//}

func (b *Bot) sendMediaErrorMessage(chatID int64) {
	b.SendMessage(chatID, "❌ Файлы, фото/видео и другие медиа <b>не принимаются</b>.")
}

// TODO передавать какой-то объект для определения нужности markup
func (b *Bot) SendMarkupMessage(chatID int64, text string) {
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

	if _, err := b.bot.Send(msg); err != nil {
		log.Fatalln(err)
	}
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := b.bot.Send(msg); err != nil {
		log.Fatalln(err)
	}
}

func (b *Bot) SendPhoto(chatID int64, path string) {
	//photo := tgbotapi.PhotoConfig{
	//	BaseFile: tgbotapi.BaseFile{
	//		BaseChat: tgbotapi.BaseChat{
	//			ChatID: chatID,
	//		},
	//		File: tgbotapi.FilePath(path),
	//	},
	//	Thumb:           nil,
	//	Caption:         caption,
	//	ParseMode:       tgbotapi.ModeHTML,
	//	CaptionEntities: nil,
	//}
	//if _, err := b.bot.Send(photo); err != nil {
	//	log.Fatalln(err)
	//}
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d&photo=%s", b.TelegramApiToken, chatID, path)
	_, err := http.Get(requestURL)
	if err != nil {
		log.Fatalln(err)
	}

}
