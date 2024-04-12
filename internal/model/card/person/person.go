package person

import (
	"fmt"
)

type PersonCard struct {
	ID                  string
	Fio                 string
	City                string
	Organization        string
	JobTitle            string
	ExpertCompetencies  string
	PossibleCooperation string
	Contacts            string
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

func ToDomain(card *PersonCard) string {
	domainCard := fmt.Sprintf(personCardTemplate,
		card.Fio,
		card.City,
		card.Organization,
		card.JobTitle,
		card.ExpertCompetencies,
		card.PossibleCooperation,
		card.Contacts,
	)

	return domainCard
}
