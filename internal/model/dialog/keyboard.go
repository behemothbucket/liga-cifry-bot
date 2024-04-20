package dialog

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	PersonKeyboard       = CreateSearchMenu("personal_cards")
	OrganizationKeyboard = CreateSearchMenu("organization_cards")
	MainKeyboard         = CreateMainMenu()
	CardKeyboard         = CreateCardMenu()
	CancelKeyboard       = CreateCancelMenu()
)

func CreateSearchMenu(searchScreen string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	buttons := BtnCriterions[searchScreen]

	for i := 0; i < len(buttons); i += 2 {
		var row []tgbotapi.KeyboardButton
		if len(rows) == 0 && i+1 < len(buttons) {
			row = append(row, tgbotapi.NewKeyboardButton(buttons[i][0]))
			row = append(row, tgbotapi.NewKeyboardButton(buttons[i+1][0]))
			row = append(row, tgbotapi.NewKeyboardButton(buttons[i+2][0]))
			i++
		} else {
			if i < len(buttons) {
				row = append(row, tgbotapi.NewKeyboardButton(buttons[i][0]))
			}
			if i+1 < len(buttons) {
				row = append(row, tgbotapi.NewKeyboardButton(buttons[i+1][0]))
			}
		}
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(BtnBack),
		tgbotapi.NewKeyboardButton(BtnApply),
	})

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	return keyboard
}

func CreateMainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnSearchPerson),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnSearchOrganization),
		),
	)
}

func CreateCardMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnMenu),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func CreateCancelMenu() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnCancelSearch),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
