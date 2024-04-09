package personal_cards

type PersonalCard struct {
	ID                  string
	Fio                 string
	City                string
	Organization        string
	JobTitle            string
	ExpertCompetencies  string
	PossibleCooperation string
	Contacts            string
}

func (pc *PersonalCard) ToDomain() PersonalCard {
	c := PersonalCard{
		Fio:                 pc.Fio,
		City:                pc.City,
		Organization:        pc.Organization,
		ExpertCompetencies:  pc.ExpertCompetencies,
		JobTitle:            pc.JobTitle,
		PossibleCooperation: pc.PossibleCooperation,
		Contacts:            pc.Contacts,
	}

	return c
}
