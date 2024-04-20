package dialog

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TestKeyboardMarkup = CreateKeyboardMarkup("personal_cards")
	TestS              = CreateKeyboardMarkupTest()

	BtnKeyboard = map[string][]string{
		"personal_cards": {
			"ФИО",
			"Город",
			"Должность",
			"Экспертные компетенции",
			"Направления сотрудничества",
			"Контакты",
		},
		"organization_cards": {
			"Организация",
			"Структурное подразделение",
			"«Приоритет-2030»",
			"Город",
			"Членство в консорциуме",
			"consortium_membership",
			"Разработки отвечественного ПО",
			"Лабораторные площадки и НОЦ",
		},
	}
)

func CreateKeyboardMarkup(searchScreen string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton

	for _, btn := range BtnKeyboard[searchScreen] {
		row := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(btn),
		)
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.InputFieldPlaceholder = "TEST"
	keyboard.ResizeKeyboard = true
	return keyboard
}

func CreateKeyboardMarkupTest() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("TEST"),
		),
	)
}
