package main

import (
	"fmt"
	"strings"
)

var (
	searchButtons = map[string][]string{
		"user": {
			"ФИО",
			"Город",
			"Организация",
			"Должность",
			"Компетенции",
			"Направления сотрудничества",
		},
		"university": {
			"Организация",
			"Структурное подразделение",
			"Город",
			"Направления сотрудничества",
			"«Приоритет-2030»",
			"Членство в консорциуме",
			"Разработки отвечественного ПО",
			"Лабораторные площадки и НОЦ",
			"Компетенции",
		},
	}

	searchUserButton       = "🔍 Индивидульные карточки"
	searchUniversityButton = "🔍 Карточки организаций"
	backButton             = "⬅️ Назад"
	menuButton             = "📋 Меню"
	cancelSearchButton     = "❌ Отменить поиск"
	applyButton            = "🆗 Применить"
	searchButton           = "🔍 Искать"
	addCard                = ""

	toggleButtonPrefix = "✅ "
)

func hasPrefix(button string, prefix string) bool {
	return strings.Contains(button, prefix)
}

func removeSearchCriterion(criteria string) {
	if currentSearchScreen == "user" {
		delete(userSearchCriteria, criteria)
	} else {
		delete(universitySearchCriteria, criteria)
	}
}

func addSearchCriterion(criteria string) {
	if currentSearchScreen == "user" {
		userSearchCriteria[criteria] = criteria
	} else {
		universitySearchCriteria[criteria] = criteria
	}
}

func findButtonIndex(buttons []string, targetButton string) int {
	for i, button := range buttons {
		if button == targetButton {
			return i
		}
	}
	return -1
}

func toggleCriterionButton(button string) {
	index := findButtonIndex(searchButtons[currentSearchScreen], button)

	if hasPrefix(searchButtons[currentSearchScreen][index], toggleButtonPrefix) {
		uncheckedButton := strings.TrimPrefix(searchButtons[currentSearchScreen][index], toggleButtonPrefix)
		searchButtons[currentSearchScreen][index] = uncheckedButton
		removeSearchCriterion(uncheckedButton)
	} else {
		searchButtons[currentSearchScreen][index] = toggleButtonPrefix + button
		addSearchCriterion(button)
	}
}

func resetCriteriaButtons() {
	for _, searchScreen := range searchButtons {
		for i, button := range searchScreen {
			if hasPrefix(button, toggleButtonPrefix) {
				searchScreen[i] = strings.TrimPrefix(button, toggleButtonPrefix)
				removeSearchCriterion(button)
			}
		}
	}
	for k := range userSearchCriteria {
		delete(userSearchCriteria, k)
	}
	for k := range universitySearchCriteria {
		delete(universitySearchCriteria, k)
	}
}

func criterionButtonIsClicked(button string) string {
	for _, v := range searchButtons[currentSearchScreen] {
		if button == v {
			return button
		}
	}
	return ""
}

func getCriterion() string {
	var criterion string
	var criteria map[string]string

	if currentSearchScreen == "user" {
		criteria = userSearchCriteria
	} else {
		criteria = universitySearchCriteria
	}

	for _, v := range criteria {
		criterion = fmt.Sprintf("Введите критерий поиска <b>%s</b>", v)
		currentCriterion = v
	}
	return criterion
}
