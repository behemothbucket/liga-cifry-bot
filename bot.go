package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type Bot struct {
	bot              *tgbotapi.BotAPI
	TelegramApiToken string
}

func newBot() *Bot {
	token := getToken()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &Bot{bot: bot, TelegramApiToken: token}
}

func getToken() string {
	token, exists := os.LookupEnv("TELEGRAM_BOT_TOKEN")

	if !exists {
		log.Print("")
	}

	return token
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if message.From == nil {
		return
	}

	logMessage(message)

	switch {
	case message.IsCommand() && getChatType(message) == "private":
		b.handleCommand(message)
	case searchMode:
		b.sendAcceptMessage(message)
	case handleIfSubscriptionEvent(message):
	case !isValidMessageText(message) && getChatType(message) == "private":
		// b.sendMediaErrorMessage(message.Chat.ID)
		b.sendMainMenu(message)
	case getChatType(message) == "private":
		b.sendMainMenu(message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	command := message.Text
	chatID := message.Chat.ID
	botName := fmt.Sprintf("@%s", b.bot.Self.UserName)

	msg := &Message{
		chatID:      message.Chat.ID,
		text:        "",
		groupName:   message.Chat.Type,
		replyMarkup: nil,
		parseMode:   tgbotapi.ModeHTML,
	}

	switch command {
	case "/start", "/start" + botName:
		b.sendMainMenu(message)
	case "/user", "/user" + botName:
		showSearchResultsMode = true
		b.SendPhoto(chatID, "https://i.imgur.com/Gyk0eeI.png")
		msg.text = userCardElena
		msg.replyMarkup = &backToMainMenuMarkup
		b.SendMessage(msg)
		showSearchResultsMode = false
	case "/university", "/university@" + botName:
		showSearchResultsMode = true
		msg.text = organizationCard
		msg.replyMarkup = &backToMainMenuMarkup
		b.SendMessage(msg)
		showSearchResultsMode = false
	}
}

func getChatType(message *tgbotapi.Message) string {
	return message.Chat.Type
}

func (b *Bot) handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := mainMenuMarkup
	message := query.Message

	switch query.Data {
	case searchUserButton:
		currentSearchScreen = "user"
		text = searchMenuDescription
		markup = getCurrentSearchMarkup()
	case searchUniversityButton:
		currentSearchScreen = "university"
		text = searchMenuDescription
		markup = getCurrentSearchMarkup()
	case backButton:
		text = mainMenuDescription
		markup = mainMenuMarkup
	case menuButton:
		//resetCriteriaButtons() // сбрасывать кнопки и чистить критерии после найденной карточки, а не здесь
		text = mainMenuDescription
		b.sendMainMenu(message)
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		return
	case applyButton:
		var criteria map[string]string

		if currentSearchScreen == "user" {
			criteria = userSearchCriteria
		} else {
			criteria = universitySearchCriteria
		}

		if len(criteria) == 0 {
			text = "️❗️Пожалуйста, выберите хотя-бы один критерий поиска."
			markup = getCurrentSearchMarkup()
		} else {
			text = getCriterion()
			markup = getCancelMenuMarkup()
			searchMode = true
		}
	case cancelSearchButton:
		resetCriteriaButtons()
		searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		b.sendMainMenu(message)
		return
	case criterionButtonIsClicked(query.Data):
		toggleCriterionButton(query.Data)
		text = searchMenuDescription
		markup = getCurrentSearchMarkup()
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.bot.Send(msg)
}

func getCurrentSearchMarkup() tgbotapi.InlineKeyboardMarkup {
	if currentSearchScreen == "user" {
		return getSearchMenuMarkup("user")
	} else {
		return getSearchMenuMarkup("university")
	}
}
