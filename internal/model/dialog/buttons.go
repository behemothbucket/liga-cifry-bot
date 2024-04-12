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

func HandleCriterionButton(button string, se search.SearchEngine) string {
	searcScreen := se.GetSearchScreen()
	buttons := btnsCriterions[searcScreen]

	for i, expected := range buttons {
		if button == expected {
			if strings.HasPrefix(buttons[i], btnChosenPrefix) {
				uncheckedButton := strings.TrimPrefix(
					buttons[i],
					btnChosenPrefix,
				)
				buttons[i] = uncheckedButton
				se.RemoveCriterion(uncheckedButton)
			} else {
				buttons[i] = btnChosenPrefix + button
				se.AddCriterion(button)
			}
		}
	}

	return button
}

func ResetCriteriaButtons() {
	for _, searchScreen := range btnsCriterions {
		for i, btn := range searchScreen {
			if strings.HasPrefix(btn, btnChosenPrefix) {
				searchScreen[i] = strings.TrimPrefix(btn, btnChosenPrefix)
			}
		}
	}
}
