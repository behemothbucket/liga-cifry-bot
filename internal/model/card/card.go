package card

import (
	"fmt"
)

type Card struct {
	Person       *PersonCard
	Organization *OrganizationCard
}

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

type OrganizationCard struct {
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

🏛 *Организация*
%s

🤝 *Должность*
%s

📝 *Экспертные компетенции*
%s

🤝 *Направления возможного сотрудничества*
%s

📱*Контакты для связи*
%s`

// TODO подсвечивать найденный текст в карточке
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

func FormatCards(cards []PersonCard) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := ToDomain(&card)
		domainCards = append(domainCards, domainCard)
	}

	return domainCards
}
