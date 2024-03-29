package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	mainMenuDescription   = "Выберите вариант поиска"
	searchMenuDescription = "<b>Выберите критерии поиска</b>"
)

var (
	mainMenuMarkup       = getMainMenuMarkup()
	cancelMenuMarkup     = getCancelMenuMarkup()
	loadMoreMenuMarkup   = getLoadMoreMenuMarkup()
	backToMainMenuMarkup = getBackToMainMenuMarkup()
)

func (b *Bot) sendMainMenu(message *tgbotapi.Message) {
	msg := Message{
		chatID:      message.Chat.ID,
		text:        mainMenuDescription,
		groupName:   message.Chat.Type,
		replyMarkup: &mainMenuMarkup,
		parseMode:   tgbotapi.ModeHTML,
	}
	b.SendMessage(msg)
}

func getMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUserButton, searchUserButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUniversityButton, searchUniversityButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(printFirstPersonalCard, printFirstPersonalCard),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(printAllPersonalCards, printAllPersonalCards),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(printFirstOrganizationCard, printFirstOrganizationCard),
		),
	)
}

func getSearchMenuMarkup(searchScreen string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range searchButtons[searchScreen] {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btn, btn),
		)
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		tgbotapi.NewInlineKeyboardButtonData(applyButton, applyButton),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func getCancelMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(cancelSearchButton, cancelSearchButton)),
	)
}

func getBackToMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(menuButton, menuButton)),
	)
}

func getLoadMoreMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(loadMoreButton, loadMoreButton)),
	)
}
