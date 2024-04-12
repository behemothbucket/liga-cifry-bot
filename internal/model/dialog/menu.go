package dialog

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	markupMainMenu               = CreateMainMenuMarkup()
	markupSearchPersonMenu       = CreateSearchMenuMarkup("person")
	markupSearchOrganizationMenu = CreateSearchMenuMarkup("organization")
	markupCancelMenu             = CreateCancelMenuMarkup()
)

func CreateMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnSearchPerson, btnSearchPerson),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnSearchOrganization, btnSearchOrganization),
		),
	)
}

func CreateSearchMenuMarkup(mode string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range btnsCriterions[mode] {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btn, btn),
		)
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(btnBack, btnBack),
		tgbotapi.NewInlineKeyboardButtonData(btnApply, btnApply),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func CreateCancelMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnCancelSearch, btnCancelSearch),
		),
	)
}
