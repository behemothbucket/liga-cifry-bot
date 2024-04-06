package personal_cards

type PersonalCard struct {
	ID                   string
	Fio                  string
	City                 string
	Organization         string
	Job_title            string
	Expert_competencies  string
	Possible_cooperation string
	Contacts             string
}

func (pc *PersonalCard) ToDomain() PersonalCard {
	c := PersonalCard{
		Fio:                  pc.Fio,
		City:                 pc.City,
		Organization:         pc.Organization,
		Job_title:            pc.Job_title,
		Possible_cooperation: pc.Possible_cooperation,
		Contacts:             pc.Contacts,
	}

	return c
}
