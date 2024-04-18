package search

import (
	"context"
	card "telegram-bot/internal/model/card"
	"telegram-bot/internal/model/db"
)

type Search struct {
	searchScreen     string
	chosenCriterions map[string]string
	searchData       []string
	isEnabled        bool
}

// Engine Интерфейс для работы с поиском карточек.
type Engine interface {
	AddCriterion(alias, criterion string)
	RemoveCriterion(alias string)
	GetCriterions() map[string]string
	ResetCriterias()
	AddSearchData(data string)
	ResetSearchData()
	GetSearchScreen() string
	SetSearchScreen(searchScreen string)
	GetSearchData() []string
	IsEnabled() bool
	Enable()
	Disable()
	ProcessCards(ctx context.Context, storage db.UserDataStorage) ([]string, error)
}

func Init() *Search {
	return &Search{}
}

func (s *Search) ProcessCards(
	ctx context.Context,
	storage db.UserDataStorage,
) ([]string, error) {
	var criteria string
	for _, v := range s.GetCriterions() {
		criteria = v
	}

	searchScreen := s.GetSearchScreen()
	searchData := s.GetSearchData()

	if searchScreen == "personal_cards" {
		rawCards, err := storage.FindPersonCards(
			ctx,
			s.GetSearchScreen(),
			searchData,
			[]string{criteria},
		)
		if err != nil {
			return nil, err
		}
		return card.FormatCardsAndHighlightPerson(rawCards, true, searchData), nil

	} else {
		rawCards, err := storage.FindOrganizationCards(
			ctx,
			s.GetSearchScreen(),
			searchData,
			[]string{criteria},
		)
		if err != nil {
			return nil, err
		}
		return card.FormatCardsAndHighlightOrganization(rawCards, true, searchData), nil
	}
}

// GetSearchScreen Получить текущий экран поиска
func (s *Search) GetSearchScreen() string {
	return s.searchScreen
}

// SetSearchScreen Назначить текущий экран поиска
func (s *Search) SetSearchScreen(screen string) {
	s.searchScreen = screen
}

func (s *Search) IsEnabled() bool {
	return s.isEnabled
}

// Disable Выключить режим поиска.
func (s *Search) Disable() {
	s.ResetCriterias()
	s.ResetSearchData()
	s.isEnabled = false
}

// Enable Включить режим поиска
func (s *Search) Enable() {
	s.isEnabled = true
}

func (s *Search) AddCriterion(alias, criterion string) {
	if s.chosenCriterions == nil {
		s.chosenCriterions = make(map[string]string)
	}
	s.chosenCriterions[alias] = criterion
}

func (s *Search) RemoveCriterion(alias string) {
	delete(s.chosenCriterions, alias)
}

func (s *Search) AddSearchData(data string) {
	s.searchData = append(s.searchData, data)
}

func (s *Search) GetCriterions() map[string]string {
	return s.chosenCriterions
}

func (s *Search) GetSearchData() []string {
	return s.searchData
}

func (s *Search) ResetCriterias() {
	s.chosenCriterions = nil
}

func (s *Search) ResetSearchData() {
	s.searchData = nil
}
