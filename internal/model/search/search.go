package search

type Search struct {
	searchScreen     string
	chosenCriterions []string
	inEnabled        bool
}

// SearchEngine Интерфейс для работы с поиском карточек.
type SearchEngine interface {
	AddCriterion(criterion string)
	RemoveCriterion(criterion string)
	ResetSearchCriterions()
	GetSearchScreen() string
	SetSearchScreen(searchScreen string)
	Disable()
	Enable()
	IsEnabled() bool
	GetCriterions() []string
}

func Init() SearchEngine {
	return &Search{}
}

// Get Получить текущий экран поиска
func (s *Search) GetSearchScreen() string {
	return s.searchScreen
}

// Set Назначить текущий экран поиска
func (s *Search) SetSearchScreen(screen string) {
	s.searchScreen = screen
}

func (s *Search) IsEnabled() bool {
	return s.inEnabled
}

// Disable Выключить режим поиска.
func (s *Search) Disable() {
	s.ResetSearchCriterions()
	s.inEnabled = false
}

// Включить режим поиска
func (s *Search) Enable() {
	s.inEnabled = true
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

func (s *Search) ResetSearchCriterions() {
	s.chosenCriterions = nil
}
