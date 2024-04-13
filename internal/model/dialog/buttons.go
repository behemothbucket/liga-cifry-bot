package dialog

import (
	"strings"
	"telegram-bot/internal/model/search"
)

// Кнопки.
var (
	BtnCriterions = map[string]map[string][]string{
		"person_cards": {
			"ФИО":         {"fio"},
			"Город":       {"city"},
			"Организация": {"organization"},
			"Должность":   {"job_title"},
			"Компетенции": {"competencies"},
			"Направления сотрудничества": {"possible_cooperation"},
		},
		"organization_cards": {
			"Структурное подразделение":     {"structural_subdivision"},
			"«Приоритет-2030»":              {"priority_2030"},
			"Членство в консорциуме":        {"consortium_membership"},
			"Разработки отвечественного ПО": {"responsible_software_development"},
			"Лабораторные площадки и НОЦ":   {"laboratories_centers"},
			"Компетенции":                   {"competencies"},
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

	for btn := range buttons {
		expected := BtnCriterions[searchScreen][btn][0]
		if button == expected {
			if strings.HasPrefix(buttons[btn][0], btnChosenPrefix) {
				uncheckedButton := strings.TrimPrefix(
					buttons[btn][0],
					btnChosenPrefix,
				)
				BtnCriterions[searchScreen][btn][0] = uncheckedButton
				se.RemoveCriterion(uncheckedButton)
			} else {
				BtnCriterions[searchScreen][btn][0] = btnChosenPrefix + button
				se.AddCriterion(buttons[btn][0])
			}
		}
	}

	return button
}

func ResetCriteriaButtons() {
	for searchScreen, buttons := range BtnCriterions {
		for btn := range buttons {
			if strings.HasPrefix(BtnCriterions[searchScreen][btn][0], btnChosenPrefix) {
				BtnCriterions[searchScreen][btn][0] = strings.TrimPrefix(
					BtnCriterions[searchScreen][btn][0],
					btnChosenPrefix,
				)
			}
		}
	}
}
