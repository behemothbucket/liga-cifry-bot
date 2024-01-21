package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	enabledKeyboard = true
)

var authorizationKeyboardButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🚪 Войти в группу Лига Цифры"),
	),
)

var menuKeyboardButton = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📜 Меню"),
	),
)
