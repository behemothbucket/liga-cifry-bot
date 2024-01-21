package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

var (
	searchButtons = map[string][]string{
		"user": {
			"ФИО",
			"Город",
			"ВУЗ",
			"Должность",
			"Компетенции",
			"Запросы на помощь/консультацию",
			"Сотрудничество",
		},
		"university": {
			"ВУЗ",
			"Приоритет 2030",
			"Кампус мирового уровня",
			"Наличие собственных разработок...",
			"Наличие лабораторных площадок...",
			"Компетенции",
			"Запросы на помощь/консультацию",
			"Сотрудничество",
		},
	}

	searchUserButton       = "Поиск карточки пользователя"
	searchUniversityButton = "Поиск карточки университета"
	backButton             = "⬅️ Назад"
	cancelButton           = "❌ Отмена"
	applyButton            = "📝 Применить"
	searchButton           = "🔍 Искать"
	miroButton             = "Miro"
)

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := getMainMenuMarkup()
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
		markup = getMainMenuMarkup()
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
	} else if query.Data == applyButton {
		text = getCriteria()
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(cancelButton, cancelButton)),
		)
		searchMode = true
	} else if query.Data == cancelButton {
		SendMenu(message.Chat.ID)
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
		searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		bot.Send(callbackCfg)
		return
	} else if criteriaButtonIsClicked(query.Data) {
		toggleButtonCheck(query.Data)
		text = searchMenuDescription
		if currentSearchScreen == "user" {
			markup = getUserSearchMenuMarkup()
		} else {
			markup = getUniversitySearchMenuMarkup()
		}
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func hasPrefix(button string, prefix string) bool {
	return strings.Contains(button, prefix)
}

func removeKey(criteria string) {
	delete(searchCriterias, criteria)
}

func findButtonIndex(buttons []string, targetButton string) int {
	for i, button := range buttons {
		if button == targetButton {
			return i
		}
	}
	return -1
}

func toggleButtonCheck(button string) {
	prefix := "✅ "

	index := findButtonIndex(searchButtons[currentSearchScreen], button)

	if hasPrefix(searchButtons[currentSearchScreen][index], prefix) {
		key := strings.TrimPrefix(searchButtons[currentSearchScreen][index], prefix)
		searchButtons[currentSearchScreen][index] = key
		removeKey(key)
	} else {
		searchButtons[currentSearchScreen][index] = prefix + button
		searchCriterias[button] = button
	}
}

func criteriaButtonIsClicked(button string) bool {
	flag := false

	for _, v := range searchButtons[currentSearchScreen] {
		if button == v {
			flag = true
			break
		}
	}

	return flag
}

func getCriteria() string {
	val := ""
	for _, v := range searchCriterias {
		val = fmt.Sprintf("Введите критерий поиска <b>%s</b>", v)
		currentCriteria = v
	}
	return val
}
