package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	id        int64
	userName  string
	firstName string
	lastName  string
	isBot     bool
}

type Message struct {
	chatID      int64
	text        string
	groupName   string
	replyMarkup *tgbotapi.InlineKeyboardMarkup
	parseMode   string
}

func isValidMessageText(message *tgbotapi.Message) bool {
	var valid bool

	if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
		valid = true
	}

	return valid
}

func handleIfSubscriptionEvent(message *tgbotapi.Message) bool {
	var event bool

	if len(message.NewChatMembers) != 0 {
		go handleNewChatMembersEvent(&message.NewChatMembers[0])
		event = true
	}
	if message.LeftChatMember != nil {
		go handleLeftChatMemberEvent(message.LeftChatMember)
		event = true
	}

	return event
}

func handleNewChatMembersEvent(user *tgbotapi.User) {
	AddUserSql(user.ID, user.UserName, user.FirstName, user.LastName, user.IsBot)
}

func handleLeftChatMemberEvent(user *tgbotapi.User) {
	DeleteUserSql(user.ID, user.UserName)
}

func logMessage(message *tgbotapi.Message) {
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

	log.Printf("https://t.me/%s [ID:%d] (%s%s) send message '%s' to chat [chatID:%d, group:%s]", userName, userID, firstName, lastName, text, chatID, groupName)
}

func (b *Bot) sendAcceptMessage(message *tgbotapi.Message) {
	msg := &Message{
		chatID:      message.Chat.ID,
		text:        "<b>Ответ принят</b>\nЯ пока что в разработке...",
		groupName:   message.Chat.Type,
		replyMarkup: &cancelMenuMarkup,
		parseMode:   tgbotapi.ModeHTML,
	}
	b.SendMessage(msg)
}

//func (b *Bot) sendSearchFinalMessage(chatID int64) {
//	b.SendMessage(chatID, "✅ Вы заполнили все критерии")
//	showSearchResultsMode = false
//	b.sendMainMenu(chatID)
//}

// func (b *Bot) sendMediaErrorMessage(chatID int64) {
// 	b.SendMessage(chatID, "❌ Файлы, фото/видео и другие медиа <b>не принимаются</b>")
// }

// TODO передавать структуру для определения нужности markup
//func (b *Bot) SendMarkupMessage(chatID int64, text string) {
//	msg := tgbotapi.NewMessage(chatID, text)
//	msg.ParseMode = tgbotapi.ModeHTML
//
//	if searchMode {
//		msg.ReplyMarkup = cancelMenuMarkup
//	}
//
//	if !searchMode && !showSearchResultsMode {
//		msg.ReplyMarkup = mainMenuMarkup
//	}
//
//	if showSearchResultsMode {
//		msg.ReplyMarkup = backToMainMenuMarkup
//	}
//
//	if _, err := b.bot.Send(msg); err != nil {
//		log.Fatalln(err)
//	}
//}

func (b *Bot) SendMessage(message *Message) {
	msg := tgbotapi.NewMessage(message.chatID, message.text)
	msg.ParseMode = message.parseMode
	if message.replyMarkup != nil {
		msg.ReplyMarkup = message.replyMarkup
	}
	if _, err := b.bot.Send(msg); err != nil {
		log.Fatalln(err)
	}
}

//func (b *Bot) SendMessage(chatID int64, text string) {
//	msg := tgbotapi.NewMessage(chatID, text)
//	msg.ParseMode = tgbotapi.ModeHTML
//	if _, err := b.bot.Send(msg); err != nil {
//		log.Fatalln(err)
//	}
//}

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
