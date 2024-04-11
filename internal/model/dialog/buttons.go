package dialog

import (
	"strings"
	"telegram-bot/internal/model/search"
)

// Кнопки.
var (
	btnsCriterions = map[string][]string{
		"person": {
			"ФИО",
			"Город",
			"Организация",
			"Должность",
			"Экспертные компетенции",
			"Направления сотрудничества",
		},
		"organization": {
			"Организация",
			"Структурное подразделение",
			"«Приоритет-2030»",
			"Город",
			"Членство в консорциуме",
			"Разработки отвечественного ПО",
			"Лабораторные площадки и НОЦ",
			"Компетенции",
		},
	}

	btnSearchPerson       = "🔍 Поиск индивидуальных карточек"
	btnSearchOrganization = "🔍 Поиск карточек организаций"
	btnBack               = "⬅️ Назад"
	// btnMenu               = "↩️ Меню"
	btnCancelSearch = "❌ Отменить поиск"
	btnApply        = "✅ Применить"
	// btnSearch             = "🔍 Искать"
	// btnLoadMore           = "⏬ Загрузить еще 5"
	btnChosenPrefix = "☑️ "
)

func IsCriterionButton(button string, mode string) string {
	for _, v := range btnsCriterions[mode] {
		if button == v {
			return button
		}
	}

	return ""
}

func findButtonIndex(buttons []string, targetButton string) int {
	for i, button := range buttons {
		if button == targetButton {
			return i
		}
	}
	return -1
}

func toggleCriterionButton(button string, se search.SearchEngine) {
	mode := se.GetMode()
	buttons := btnsCriterions[mode]
	index := findButtonIndex(buttons, button)

	if hasPrefix(buttons[index], btnChosenPrefix) {
		uncheckedButton := strings.TrimPrefix(
			buttons[index],
			btnChosenPrefix,
		)
		buttons[index] = uncheckedButton
		se.RemoveCriterion(uncheckedButton)
	} else {
		buttons[index] = btnChosenPrefix + button
		se.AddCriterion(button)
	}
}

func hasPrefix(button string, prefix string) bool {
	return strings.Contains(button, prefix)
}
