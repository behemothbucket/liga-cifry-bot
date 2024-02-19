package main

import (
	"context"
	"log"
	"reflect"
	sqlite "telegram-bot/storage/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func handleIfSubscriptionEvent(ctx context.Context, message *tgbotapi.Message) bool {
	var event bool

	s, err := sqlite.New("users.db")
	if err != nil {
		log.Panicf("can't connect to storage: %v ", err)
	}

	if err := s.Init(ctx); err != nil {
		log.Panicf("can't init storage: %v ", err)
	}

	if len(message.NewChatMembers) != 0 {
		go s.AddUser(ctx, &message.NewChatMembers[0])
		event = true
	}
	if message.LeftChatMember != nil {
		go s.DeleteUser(ctx, message.LeftChatMember)
		event = true
	}

	return event
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

	log.Printf("https://t.me/%s [ID:%d] (%s%s) send message '%s' to chat [chatID:%d, group:%s]",
		userName, userID, firstName, lastName, text, chatID, groupName)
}

func (b *Bot) sendAcceptMessage(message *tgbotapi.Message) {
	msg := Message{
		chatID:      message.Chat.ID,
		text:        "<b>Ответ принят</b>\nЯ пока что в разработке...",
		groupName:   message.Chat.Type,
		replyMarkup: &cancelMenuMarkup,
		parseMode:   tgbotapi.ModeHTML,
	}
	b.SendMessage(msg)
}

func (b *Bot) sendLoadMoreMessage(message *tgbotapi.Message) {
	msg := Message{
		chatID:      message.Chat.ID,
		text:        "<b>Загрузить больше вариантов</b>",
		groupName:   message.Chat.Type,
		replyMarkup: &loadMoreMenuMarkup,
		parseMode:   tgbotapi.ModeHTML,
	}
	b.SendMessage(msg)
}

func (b *Bot) SendMessage(message Message) {
	msg := tgbotapi.NewMessage(message.chatID, message.text)
	msg.ParseMode = message.parseMode
	if message.replyMarkup != nil {
		msg.ReplyMarkup = message.replyMarkup
	}
	if _, err := b.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
