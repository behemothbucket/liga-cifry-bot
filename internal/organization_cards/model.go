package cards

type OrganizationCard struct {
	Name                   string
	Structural_subdivision string
	City                   string
	Possible_cooperation   string
	Priority_2030          bool
	Consortium_membership  bool
	Software               string
	Laboratory_and_noc     bool
}

func (oc *OrganizationCard) ToDomain() OrganizationCard {
	c := OrganizationCard{
		Name:                   oc.Name,
		Structural_subdivision: oc.Structural_subdivision,
		City:                   oc.City,
		Possible_cooperation:   oc.Possible_cooperation,
		Priority_2030:          oc.Priority_2030,
		Consortium_membership:  oc.Consortium_membership,
		Software:               oc.Software,
		Laboratory_and_noc:     oc.Laboratory_and_noc,
	}

	return c
}
