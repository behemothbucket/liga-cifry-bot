package search

import (
	"context"
	"telegram-bot/internal/model/card/person"
	"telegram-bot/internal/model/db"
)

type Search struct {
	searchScreen     string
	chosenCriterions []string
	isEnabled        bool
}

// Engine Интерфейс для работы с поиском карточек.
type Engine interface {
	AddCriterion(criterion string)
	RemoveCriterion(criterion string)
	ResetSearchCriterias()
	GetSearchScreen() string
	SetSearchScreen(searchScreen string)
	Disable()
	Enable()
	IsEnabled() bool
	FormatCards(cards []person.PersonCard) []string
	ProcessCards(ctx context.Context, storage db.UserDataStorage) ([]string, error)
	GetCriterions() []string
}

func Init() Engine {
	return &Search{
		isEnabled: false,
	}
}

func (s *Search) ProcessCards(ctx context.Context, storage db.UserDataStorage) ([]string, error) {
	rawCards, err := storage.FindCards(ctx, s.GetCriterions()[0], s.chosenCriterions)
	if err != nil {
		return make([]string, 0), err
	}

	return s.FormatCards(rawCards), nil
}

func (s *Search) FormatCards(cards []person.PersonCard) []string {
	var domainCards []string

	for _, card := range cards {
		domainCard := person.ToDomain(&card)
		domainCards = append(domainCards, domainCard)
	}

	return domainCards
}

// GetSearchScreen Get Получить текущий экран поиска
func (s *Search) GetSearchScreen() string {
	return s.searchScreen
}

// SetSearchScreen Set Назначить текущий экран поиска
func (s *Search) SetSearchScreen(screen string) {
	s.searchScreen = screen
}

func (s *Search) IsEnabled() bool {
	return s.isEnabled
}

// Disable Выключить режим поиска.
func (s *Search) Disable() {
	s.ResetSearchCriterias()
	s.isEnabled = false
}

// Enable Включить режим поиска
func (s *Search) Enable() {
	s.isEnabled = true
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

func (s *Search) ResetSearchCriterias() {
	s.chosenCriterions = nil
}
