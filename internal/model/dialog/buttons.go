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
			{"Направления сотрудничества", "possible_cooperation"},
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
	BtnMenu               = "↩️ Меню"
	BtnCancelSearch       = "❌ Отменить поиск"
	BtnApply              = "✅ Применить"
	// btnLoadMore           = "⏬ Загрузить еще 5"
	BtnChosenPrefix      = "☑️ "
	BtnTestReplyKeyboard = "TEST"
)

func HandleCriterionButton(button string, se search.Engine) string {
	searchScreen := se.GetSearchScreen()
	buttons := BtnCriterions[searchScreen]

	for i, expected := range buttons {
		if button == expected[1] {
			if strings.HasPrefix(buttons[i][0], BtnChosenPrefix) {
				uncheckedButton := strings.TrimPrefix(
					expected[0],
					BtnChosenPrefix,
				)
				buttons[i][0] = uncheckedButton
				se.RemoveCriterion(BtnChosenPrefix + expected[0])
			} else {
				buttons[i][0] = BtnChosenPrefix + expected[0]
				se.AddCriterion(buttons[i][0], button)
			}
			break
		}
	}

	return button
}

func ResetCriteriaButtons() {
	for _, searchScreen := range BtnCriterions {
		for i, btn := range searchScreen {
			searchScreen[i][0] = strings.TrimPrefix(btn[0], BtnChosenPrefix)
		}
	}
}
