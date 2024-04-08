package personal_cards

type PersonalCard struct {
	ID                   string
	fio                  string
	city                 string
	organization         string
	job_title            string
	expert_competencies  string
	possible_cooperation string
	contacts             string
}

func (pc *PersonalCard) ToDomain() PersonalCard {
	c := PersonalCard{
		fio:                  pc.fio,
		city:                 pc.city,
		organization:         pc.organization,
		job_title:            pc.job_title,
		expert_competencies:  pc.expert_competencies,
		possible_cooperation: pc.possible_cooperation,
		contacts:             pc.contacts,
	}

	return c
}
