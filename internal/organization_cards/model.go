package cards

type OrganizationCard struct {
	name                  string
	structuralSubdivision string
	city                  string
	possibleCooperation   string
	priority2030          bool
	consortiumMembership  bool
	software              string
	laboratoryAndNoc      bool
}

func (oc *OrganizationCard) ToDomain() OrganizationCard {
	c := OrganizationCard{
		name:                  oc.name,
		structuralSubdivision: oc.structuralSubdivision,
		city:                  oc.city,
		possibleCooperation:   oc.possibleCooperation,
		priority2030:          oc.priority2030,
		consortiumMembership:  oc.consortiumMembership,
		software:              oc.software,
		laboratoryAndNoc:      oc.laboratoryAndNoc,
	}

	return c
}
