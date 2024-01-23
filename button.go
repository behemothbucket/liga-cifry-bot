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
			"Направления сотрудничества",
			"«Приоритет-2030»",
			"Кампус мирового уровня",
			"Разработки отвечественного ПО",
			"Лабораторные площадки и НОЦ",
			"Компетенции",
		},
	}

	searchUserButton       = "🔍 Индивидульные карточки"
	searchUniversityButton = "🔍 Карточки организаций"
	backButton             = "⬅️ Назад"
	cancelButton           = "❌ Отменить поиск"
	applyButton            = "🆗 Применить"
	searchButton           = "🔍 Искать"
	addCard                = ""

	toggleButtonPrefix = "✅ "
)

func hasPrefix(button string, prefix string) bool {
	return strings.Contains(button, prefix)
}

func removeSearchCriteria(criteria string) {
	delete(searchCriterias, criteria)
}

func addSearchCriteria(criteria string) {
	searchCriterias[criteria] = criteria
}

func removeAllSearchCriterias() {
	for k := range searchCriterias {
		delete(searchCriterias, k)
	}
}

func removeCriteriaByPrefix(screen []string, prefix string) {
	for i, _ := range screen {
		if hasPrefix(screen[i], prefix) {
			key := strings.TrimPrefix(screen[i], prefix)
			screen[i] = key
			removeSearchCriteria(key)
		}
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

func toggleCriteriaButton(button string) {
	index := findButtonIndex(searchButtons[currentSearchScreen], button)

	if hasPrefix(searchButtons[currentSearchScreen][index], toggleButtonPrefix) {
		removedPrefix := strings.TrimPrefix(searchButtons[currentSearchScreen][index], toggleButtonPrefix)
		searchButtons[currentSearchScreen][index] = removedPrefix
		removeSearchCriteria(removedPrefix)
	} else {
		searchButtons[currentSearchScreen][index] = toggleButtonPrefix + button
		addSearchCriteria(button)
	}
}

func resetCriteriaButtons() {
	removeCriteriaByPrefix(searchButtons[currentSearchScreen], toggleButtonPrefix)
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
