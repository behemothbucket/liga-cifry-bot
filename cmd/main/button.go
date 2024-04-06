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

	searchUserButton           = "🔍 Индивидульные карточки"
	searchUniversityButton     = "🔍 Карточки организаций"
	backButton                 = "⬅️ Назад"
	menuButton                 = "↩️ Меню"
	cancelSearchButton         = "❌ Отменить поиск"
	applyButton                = "🆗 Применить"
	searchButton               = "🔍 Искать"
	printFirstPersonalCard     = "⚠️Персональная карточка⚠️"
	printAllPersonalCards      = "⚠️Все персональные карточки⚠️"
	printFirstOrganizationCard = "⚠️Карточка организации⚠️"
	loadMoreButton             = "⏬ Загрузить еще 5"
	// addCard                    = ""

	toggleButtonPrefix = "✅ "
)

func hasPrefix(button string, prefix string) bool {
	return strings.Contains(button, prefix)
}

func (b *Bot) removeSearchCriterion(criteria string) {
	if b.currentSearchScreen == "user" {
		delete(b.userSearchCriteria, criteria)
	} else {
		delete(b.universitySearchCriteria, criteria)
	}
}

func (b *Bot) addSearchCriterion(criteria string) {
	if b.currentSearchScreen == "user" {
		b.userSearchCriteria[criteria] = criteria
	} else {
		b.universitySearchCriteria[criteria] = criteria
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

func (b *Bot) toggleCriterionButton(button string) {
	index := findButtonIndex(searchButtons[b.currentSearchScreen], button)

	if hasPrefix(searchButtons[b.currentSearchScreen][index], toggleButtonPrefix) {
		uncheckedButton := strings.TrimPrefix(
			searchButtons[b.currentSearchScreen][index],
			toggleButtonPrefix,
		)
		searchButtons[b.currentSearchScreen][index] = uncheckedButton
		b.removeSearchCriterion(uncheckedButton)
	} else {
		searchButtons[b.currentSearchScreen][index] = toggleButtonPrefix + button
		b.addSearchCriterion(button)
	}
}

func (b *Bot) resetCriteriaButtons() {
	for _, searchScreen := range searchButtons {
		for i, button := range searchScreen {
			if hasPrefix(button, toggleButtonPrefix) {
				searchScreen[i] = strings.TrimPrefix(button, toggleButtonPrefix)
				b.removeSearchCriterion(button)
			}
		}
	}
	for k := range b.userSearchCriteria {
		delete(b.userSearchCriteria, k)
	}
	for k := range b.universitySearchCriteria {
		delete(b.universitySearchCriteria, k)
	}
}

func (b *Bot) criterionButtonIsClicked(button string) string {
	for _, v := range searchButtons[b.currentSearchScreen] {
		if button == v {
			return button
		}
	}

	return ""
}

func (b *Bot) getCriterion() string {
	var criterion string
	var criteria map[string]string

	if b.currentSearchScreen == "user" {
		criteria = b.userSearchCriteria
	} else {
		criteria = b.universitySearchCriteria
	}

	for _, v := range criteria {
		criterion = fmt.Sprintf("Введите критерий поиска <b>%s</b>", v)
		// currentCriterion = v
	}
	return criterion
}
