package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"os"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocket.Client
}

func NewBot() *Bot {
	var err error

	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Panic(err)
	}

	return &Bot{bot: bot}
}

func getToken() string {
	token, exists := os.LookupEnv("TOKEN")

	if !exists {
		log.Print("Токен не обнаружен")
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
	default:
		b.handleMessage(update)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	if update.Message.From == nil {
		return
	}

	user := update.Message.From
	text := update.Message.Text

	logMessage(user, text)

	var err error

	if isValidMessageText(text) {
		if isAuthorized {
			if update.Message.IsCommand() {
				b.handleCommand(update)
			} else if searchMode {
				if len(searchCriterias) == 0 {
					searchMode = false
					err = b.sendSearchFinalMessage(update.Message.Chat.ID)
				} else {
					err = b.sendAcceptMessage(update.Message.Chat.ID)
					delete(searchCriterias, currentCriteria)
				}
			} else {
				err = b.sendMainMenu(update.Message.Chat.ID)
			}
		} else if text == "В клюве" {
			isAuthorized = true
			enabledInlineKeyboard = true
			enabledKeyboard = false
			mainMenuMarkup = getMainMenuMarkup()
			err = b.sendMainMenu(update.Message.Chat.ID)
		} else {
			err = b.sendAuthorizationMessage(update.Message.Chat.ID)
		}
	} else {
		err = b.sendMediaErrorMessage(update.Message.Chat.ID)
		err = b.sendMainMenu(update.Message.Chat.ID)
	}

	if err != nil {
		log.Println("Ошибка:", err.Error())
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	command := update.Message.Text

	switch command {
	case "/start":
		if isAuthorized {
			b.sendMainMenu(update.Message.Chat.ID)
		} else {
			b.sendAuthorizationMessage(update.Message.Chat.ID)
		}
		break
	case "/login":
		if isAuthorized {
			b.sendMessage(update.Message.Chat.ID, "✅ Вы уже авторизованы ")
			b.sendMainMenu(update.Message.Chat.ID)
		} else {
			b.sendAuthorizationMessage(update.Message.Chat.ID)
		}
		break
	case "/invite":
		b.sendMessage(update.Message.Chat.ID, "🏗 В разработке...")
		b.sendMainMenu(update.Message.Chat.ID)
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
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
	} else if query.Data == applyButton {
		if len(searchCriterias) == 0 {
			text = "️❗️Пожалуйста, выберите хотя-бы один критерий поиска"
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
