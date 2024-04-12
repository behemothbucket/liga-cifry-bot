package dialog

import (
	"strings"
	"telegram-bot/internal/model/search"
)

// Кнопки.
var (
	BtnCriterions = map[string][]string{
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
	BtnSearchPerson       = "🔍 Поиск индивидуальных карточек"
	BtnSearchOrganization = "🔍 Поиск карточек организаций"
	BtnBack               = "⬅️ Назад"
	// btnMenu               = "↩️ Меню"
	BtnCancelSearch = "❌ Отменить поиск"
	BtnApply        = "✅ Применить"
	// btnSearch             = "🔍 Искать"
	// btnLoadMore           = "⏬ Загрузить еще 5"
	btnChosenPrefix = "☑️ "
)

func HandleCriterionButton(button string, se search.Engine) string {
	searchScreen := se.GetSearchScreen()
	buttons := BtnCriterions[searchScreen]

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
				se.AddCriterion(buttons[i])
			}
		}
	}

	return button
}

func ResetCriteriaButtons() {
	for _, searchScreen := range BtnCriterions {
		for i, btn := range searchScreen {
			if strings.HasPrefix(btn, btnChosenPrefix) {
				searchScreen[i] = strings.TrimPrefix(btn, btnChosenPrefix)
			}
		}
	}
}
