package person

import (
	"fmt"
	"telegram-bot/internal/helpers/markdown"
)

type PersonCard struct {
	ID                   string
	Fio                  string
	City                 string
	Organization         string
	Job_title            string
	Expert_competencies  string
	Possible_cooperation string
	Contacts             string
}

const personCardTemplate = `
🧑‍💼*ФИО*
%s

📍*Город*
%s

🏛*Организация*
%s

🤝*Должность*
%s

📝*Экспертные компетенции*
%s

🤝*Направления возможного сотрудничества*
%s

📱*Контакты для связи*
%s`

func MarkupCard(card *PersonCard) string {
	formattedText := fmt.Sprintf(personCardTemplate,
		card.Fio,
		card.City,
		card.Organization,
		card.Job_title,
		card.Expert_competencies,
		card.Possible_cooperation,
		card.Contacts,
	)

	return markdown.EscapeForMarkdown(formattedText)
}
