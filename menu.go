package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	mainMenuDescription   = "Выберите вариант поиска"
	searchMenuDescription = "<b>Выберите критерии поиска</b>"
)

var (
	mainMenuMarkup       = getMainMenuMarkup()
	cancelMenuMarkup     = getCancelMenuMarkup()
	backToMainMenuMarkup = getBackToMainMenuMarkup()
)

func (b *Bot) sendMainMenu(chatId int64) {
	b.SendMarkupMessage(chatId, mainMenuDescription)
}

func getMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUserButton, searchUserButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUniversityButton, searchUniversityButton),
		),
	)
}

func getSearchMenuMarkup(searchType string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range searchButtons[searchType] {
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

func getUserSearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("user")
}

func getCancelMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(cancelButton, cancelButton)),
	)
}

func getBackToMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(backToMainMenuButton, backToMainMenuButton)),
	)
}

func getUniversitySearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("university")
}
