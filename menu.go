package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	mainMenuDescription   = "<b>Меню</b>\n\nТекст Текст Текст\n\nТекст Текст Текст\n\nТекст Текст Текст"
	searchMenuDescription = "<b>Выберите критерии поиска:</b>"
)

var (
	enabledInlineKeyboard = false
)

func SendMenu(chatId int64) error {
	return sendFormattedMessage(chatId, mainMenuDescription)
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
			tgbotapi.NewInlineKeyboardButtonURL(miroButton, "https://miro.com/app/board/uXjVN5NbjoM=/"),
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
		tgbotapi.NewInlineKeyboardButtonData(applyButton, applyButton),
	})

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func getUserSearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("user")
}

func getUniversitySearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("university")
}
