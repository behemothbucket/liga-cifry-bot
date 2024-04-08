package cards

type OrganizationCard struct {
	name                   string
	structural_subdivision string
	city                   string
	possible_cooperation   string
	priority_2030          bool
	consortium_membership  bool
	software               string
	laboratory_and_noc     bool
}

func (oc *OrganizationCard) ToDomain() OrganizationCard {
	c := OrganizationCard{
		name:                   oc.name,
		structural_subdivision: oc.structural_subdivision,
		city:                   oc.city,
		possible_cooperation:   oc.possible_cooperation,
		priority_2030:          oc.priority_2030,
		consortium_membership:  oc.consortium_membership,
		software:               oc.software,
		laboratory_and_noc:     oc.laboratory_and_noc,
	}

	return c
}
