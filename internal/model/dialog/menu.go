package dialog

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	MarkupMainMenu               = CreateMainMenuMarkup()
	MarkupSearchPersonMenu       = CreateSearchMenuMarkup("personal_cards")
	MarkupSearchOrganizationMenu = CreateSearchMenuMarkup("organization_cards")
	MarkupCancelMenu             = CreateCancelMenuMarkup()
	MarkupCardMenu               = CreateCardMenuMarkup()
)

func CreateMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnSearchPerson, BtnSearchPerson),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnSearchOrganization, BtnSearchOrganization),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnTestReplyKeyboard, BtnTestReplyKeyboard),
		),
	)
}

func CreateSearchMenuMarkup(searchScreen string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range BtnCriterions[searchScreen] {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btn[0], btn[1]),
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

func CreateCardMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtnMenu, BtnMenu),
		),
	)
}
