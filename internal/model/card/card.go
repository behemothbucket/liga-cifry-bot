package card

import (
	"fmt"
	"strings"
)

const personCardTemplate = `
🧑‍💼<b>ФИО</b>
%s

📍<b>Город</b>
%s

🏛 <b>Организация</b>
%s

🤝 <b>Должность</b>
%s

📝 <b>Экспертные компетенции</b>
%s

🤝 <b>Направления возможного сотрудничества</b>
%s

📱<b>Контакты для связи</b>
%s`

const organizationCardTemplate = `
🏛 <b>Организация</b>
%s

🏢 <b>Структурное подразделение</b>

📍<b>Город</b>
%s

🤝 <b>Направления возможного сотрудничества</b>
%s

🌐 <b>Членство в консорциуме</b>
%s

🚀 <b>«Приоритет-2030»</b>
%s

⚛️ <b>Разработки отечественного ПО</b>
%s

🔬 <b>Лабораторные площадки и НОЦ</b>
%s`

type Card interface {
	ToDomain() string
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
	ID                   string
	City                 string
	Organization         string
	ConsortiumMembership string
	PossibleCooperation  string
	Priority2030         string
	Contacts             string
	Software             string
	LaboratoriesAndNOCs  string
}

func (pc PersonCard) ToDomain(field string) string {
	return fmt.Sprintf(
		personCardTemplate,
		pc.Fio,
		pc.City,
		pc.Organization,
		pc.JobTitle,
		pc.ExpertCompetencies,
		pc.PossibleCooperation,
		pc.Contacts)
}

func (oc OrganizationCard) ToDomain() string {
	return fmt.Sprintf(
		organizationCardTemplate,
		oc.Organization,
		oc.City,
		oc.PossibleCooperation,
		oc.ConsortiumMembership,
		oc.Priority2030,
		oc.Software,
		oc.LaboratoriesAndNOCs)
}

// TODO подсвечивать найденный текст в карточке
func ToDomain(card Card) string {
	var domainCard string

	switch c := card.(type) {
	case PersonCard:
		domainCard = c.ToDomain()
	case OrganizationCard:
		domainCard = c.ToDomain()
	}

	return domainCard
}

func FormatPersonCards(cards []PersonCard) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := ToDomain(card)
		domainCards = append(domainCards, domainCard)
	}

	return domainCards
}

func FormatOrganizationCards(cards []OrganizationCard) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := ToDomain(card)
		domainCards = append(domainCards, domainCard)
	}

	return domainCards
}

func FormatCardsAndHighlightOrganization(
	cards []OrganizationCard,
	highlight bool,
	searchData []string,
) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := ToDomain(card)
		if highlight {
			for _, phrase := range searchData {
				if strings.Contains(strings.ToLower(domainCard), strings.ToLower(phrase)) {
					domainCard = strings.Replace(domainCard, phrase, "<code>"+phrase+"</code>", -1)
				}
			}
		}
		domainCards = append(domainCards, domainCard)
	}

	return domainCards
}

func FormatCardsAndHighlightPerson(
	cards []PersonCard,
	highlight bool,
	searchData []string,
) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := ToDomain(card)
		if highlight {
			for _, phrase := range searchData {
				if strings.Contains(strings.ToLower(domainCard), strings.ToLower(phrase)) {
					domainCard = strings.Replace(domainCard, phrase, "<code>"+phrase+"</code>", -1)
				}
			}
		}
		domainCards = append(domainCards, domainCard)
	}
	return domainCards
}
