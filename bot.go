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
		log.Print("Токен не обнаружен.")
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
		b.handleMessage(update)
		break
	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery)
		break
	}

}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	if update.Message.From == nil {
		return
	}

	chatID := update.Message.Chat.ID
	user := update.Message.From
	message := update.Message

	logMessage(user, message.Text, chatID)

	if isValidMessageText(update) && !isJoinEvent(update) {
		switch {
		case update.Message.IsCommand():
			b.handleCommand(update)
			break
		case searchMode:
			b.sendAcceptMessage(chatID)
			break
		default:
			b.sendMainMenu(chatID)
		}
	} else {
		b.sendMediaErrorMessage(update.Message.Chat.ID)
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	command := update.Message.Text
	chatType := update.Message.Chat.Type
	chatID := update.Message.Chat.ID
	botName := fmt.Sprintf("@%s", b.bot.Self.UserName)

	if chatType == "group" || chatType == "supergroup" || chatType == "channel" {
		log.Printf("Обнаружен чат с типом %s, блокирую комманду", chatType)
		b.SendMessage(chatID, "Вы не являетесь <b>администратором</b>, чтобы использовать данную комманду.")
		return
	}

	switch command {
	case "/start", "/start" + botName:
		b.sendMainMenu(chatID)
		break
	case "/user", "/user" + botName:
		showSearchResultsMode = true
		b.SendPhoto(chatID, "https://i.imgur.com/Gyk0eeI.png")
		b.SendMarkupMessage(chatID, userCardElena)
	case "/university", "/university@" + botName:
		showSearchResultsMode = true
		b.SendMessage(chatID, organizationCard)
		b.SendMarkupMessage(chatID, competitionCard)
	}
}

func (b *Bot) handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := mainMenuMarkup
	message := query.Message

	if query.Data == searchUserButton {
		text = searchMenuDescription
		markup = getUserSearchMenuMarkup()
		currentSearchScreen = "user"
	} else if query.Data == searchUniversityButton {
		text = searchMenuDescription
		markup = getUniversitySearchMenuMarkup()
		currentSearchScreen = "university"
	} else if query.Data == backButton {
		text = mainMenuDescription
		markup = mainMenuMarkup
		removeAllSearchCriterias()
	} else if query.Data == backToMainMenuButton {
		text = mainMenuDescription
		removeAllSearchCriterias()
		b.sendMainMenu(message.Chat.ID)
		return
	} else if query.Data == applyButton {
		if len(searchCriterias) == 0 {
			text = "️❗️Пожалуйста, выберите хотя-бы один критерий поиска."
			markup = getUserSearchMenuMarkup()
		} else {
			text = getCriteria()
			searchMode = true
			cancelMenuMarkup = getCancelMenuMarkup()
			markup = cancelMenuMarkup
		}
	} else if query.Data == cancelButton {
		removeAllSearchCriterias()
		resetCriteriaButtons()
		searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		b.sendMainMenu(message.Chat.ID)
		return
	} else if criteriaButtonIsClicked(query.Data) {
		toggleCriteriaButton(query.Data)
		text = searchMenuDescription
		if currentSearchScreen == "user" {
			markup = getUserSearchMenuMarkup()
		} else {
			markup = getUniversitySearchMenuMarkup()
		}
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.bot.Send(msg)
}
