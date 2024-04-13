package dialog

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	MarkupMainMenu               = CreateMainMenuMarkup()
	MarkupSearchPersonMenu       = CreateSearchMenuMarkup("person")
	MarkupSearchOrganizationMenu = CreateSearchMenuMarkup("organization")
	MarkupCancelMenu             = CreateCancelMenuMarkup()
)

func CreateMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnSearchPerson, BtnSearchPerson),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnSearchOrganization, BtnSearchOrganization),
		),
	)
}

func CreateSearchMenuMarkup(searchScreen string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for key, value := range BtnCriterions[searchScreen] {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(key, value[0]),
		)
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(BtnBack, BtnBack),
		tgbotapi.NewInlineKeyboardButtonData(BtnApply, BtnApply),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func CreateCancelMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnCancelSearch, BtnCancelSearch),
		),
	)
}
