package search

type Search struct {
	mode             string
	chosenCriterions []string
}

// SearchEngine Интерфейс для работы с поиском карточек.
type SearchEngine interface {
	AddCriterion(criterion string)
	RemoveCriterion(criterion string)
	ResetCriterions()
	GetMode() string
	SetMode(searchType string)
	Disable()
	GetCriterions() []string
}

func Init() SearchEngine {
	return &Search{}
}

// Get Получить текущий экран поиска
func (s *Search) GetMode() string {
	return s.mode
}

// Set Назначить текущий экран поиска
func (s *Search) SetMode(mode string) {
	s.mode = mode
}

// Disable Выключить режим поиска.
func (s *Search) Disable() {
	s.mode = ""
	s.ResetCriterions()
}

func (s *Search) AddCriterion(criterion string) {
	s.chosenCriterions = append(s.chosenCriterions, criterion)
}

func (s *Search) RemoveCriterion(criterion string) {
	for i, _criterion := range s.chosenCriterions {
		if _criterion == criterion {
			s.chosenCriterions = append(s.chosenCriterions[:i], s.chosenCriterions[i+1:]...)
			break
		}
	}
}

func (s *Search) GetCriterions() []string {
	return s.chosenCriterions
}

func (s *Search) ResetCriterions() {
	s.chosenCriterions = nil
}
