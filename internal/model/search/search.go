package search

import (
	"context"
	card "telegram-bot/internal/model/card"
	"telegram-bot/internal/model/db"
)

type Search struct {
	searchScreen     string
	chosenCriterions []string
	searchData       []string
	isEnabled        bool
}

// Engine Интерфейс для работы с поиском карточек.
type Engine interface {
	AddCriterion(criterion string)
	AddSearchData(data string)
	RemoveCriterion(criterion string)
	ResetCriterias()
	ResetSearchData()
	GetSearchScreen() string
	SetSearchScreen(searchScreen string)
	GetCriterions() []string
	GetSearchData() []string
	IsEnabled() bool
	Enable()
	Disable()
	ProcessCards(ctx context.Context, storage db.UserDataStorage) ([]string, error)
}

func Init() Engine {
	return &Search{
		isEnabled: false,
	}
}

func (s *Search) ProcessCards(
	ctx context.Context,
	storage db.UserDataStorage,
) ([]string, error) {
	// TODO передавать введенные пользователем данные
	rawCards, err := storage.FindCards(
		ctx,
		s.GetSearchScreen(),
		s.GetSearchData(),
		s.GetCriterions(),
	)
	if err != nil {
		return make([]string, 0), err
	}

	return card.FormatCards(rawCards), nil
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
	s.ResetCriterias()
	s.ResetSearchData()
	s.isEnabled = false
}

// Enable Включить режим поиска
func (s *Search) Enable() {
	s.isEnabled = true
}

func (s *Search) AddCriterion(criterion string) {
	s.chosenCriterions = append(s.chosenCriterions, criterion)
}

func (s *Search) AddSearchData(data string) {
	s.searchData = append(s.searchData, data)
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

func (s *Search) GetSearchData() []string {
	return s.searchData
}

func (s *Search) ResetCriterias() {
	s.chosenCriterions = nil
}

func (s *Search) ResetSearchData() {
	s.searchData = nil
}
