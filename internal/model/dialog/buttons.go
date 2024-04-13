package dialog

import (
	"strings"
	"telegram-bot/internal/model/search"
)

// Кнопки.
var (
	BtnCriterions = map[string][][]string{
		"personal_cards": {
			{"ФИО", "fio"},
			{"Город", "city"},
			{"Организация", "organization"},
			{"Должность", "job_title"},
			{"Экспертные компетенции", "expert_competencies"},
			{"Направления сотрудничества", "possible_cooperations"},
			{"Контакты", "contacts"},
		},
		"organization_cards": {
			{"Организация", "name"},
			{"Структурное подразделение", "structural_subdivision"},
			{"«Приоритет-2030»", "priority_2030"},
			{"Город", "city"},
			{"Членство в консорциуме", "consortium_membership"},
			{"Разработки отвечественного ПО", "software"},
			{"Лабораторные площадки и НОЦ", "laboratory_and_noc"},
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
		if button == expected[0] {
			if strings.HasPrefix(buttons[i][0], btnChosenPrefix) {
				uncheckedButton := strings.TrimPrefix(
					buttons[i][0],
					btnChosenPrefix,
				)
				buttons[i][0] = uncheckedButton
				se.RemoveCriterion(uncheckedButton)
			} else {
				buttons[i][0] = btnChosenPrefix + button
				se.AddCriterion(button)
			}
		}
	}

	return button
}

func ResetCriteriaButtons() {
	for _, searchScreen := range BtnCriterions {
		for i, btn := range searchScreen {
			searchScreen[i][0] = strings.TrimPrefix(btn[0], btnChosenPrefix)
		}
	}
}
