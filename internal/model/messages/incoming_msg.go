package messages

import (
	"context"
	"fmt"
	types "telegram-bot/internal/model/bottypes"

	"github.com/opentracing/opentracing-go"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Область "Константы и переменные": начало.

const (
	txtStart           = "Привет, *%v*. Я помогу найти карточку компетенций или карточку организации.\b Выберите действие."
	txtUnknownCommand  = "К сожалению, данная команда мне неизвестна. Для начала работы введите /start"
	txtReportError     = "Не удалось получить данные."
	txtReportWait      = "Ищу 🔎\nПожалуйста, подождите..."
	txtCriterionChoose = "Выберите критерии поиска для поиска , а затем нажмите **Применить**."
	txtCatEmpty        = "Не выбрано ни одного критерия поиска. Сначала выберите хотя-бы один критерий."
	txtHelp            = "Я - бот, помогающий искать карточки компетенций и организаций для проекта **Лига Цифры**. Для начала работы введите /start"
)

// Команды стартовых действий.
var btnStart = []types.TgRowButtons{
	{
		types.TgInlineButton{
			DisplayName: "🔍 Индивидульные карточки",
			Value:       "/search_personal_card",
		},
		types.TgInlineButton{
			DisplayName: "🔍 Карточки организаций",
			Value:       "/search_organization_card",
		},
		types.TgInlineButton{
			DisplayName: "⚠️Все персональные карточки⚠️",
			Value:       "/all_personal_cards",
		},
	},
}

// Область "Константы и переменные": конец.

// Область "Внешний интерфейс": начало.

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, userID int64) error
	ShowInlineButtons(text string, buttons []types.TgRowButtons, userID int64) error
}

// UserDataStorage Интерфейс для работы с хранилищем данных.
type UserDataStorage interface {
	JoinGroup(
		ctx context.Context,
		id int64,
		userName string,
		firstName string,
		lastName string,
		isBot bool,
	) error
	LeaveGroup(ctx context.Context, u *tgbotapi.User) error
	CheckIfUserExist(ctx context.Context, userID int64) (bool, error)
	ShowAllPersonalCards(
		ctx context.Context,
	) (pc []PersonalCard, err error)
}

// SearchEngine Интерфейс для работы с поиском карточек.
type SearchEngine interface {
	RemoveCriterion(criteria string)
	AddCriterion(criteria string)
	ToggleCriterion(button types.TgInlineButton)
	ResetCriteriaButtons()
}

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

// Model Модель бота (клиент, хранилище, последние команды пользователя)
type Model struct {
	ctx      context.Context
	tgClient MessageSender   // Клиент.
	storage  UserDataStorage // Хранилище пользовательской информации.
	// search          SearchEngine     // Поиск.
	lastUserCommand map[int64]string // Последняя выбранная пользователем команда.
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей и параметров поиска.
func New(
	ctx context.Context,
	tgClient MessageSender,
	storage UserDataStorage,
	// search SearchEngine,
) *Model {
	return &Model{
		ctx:      ctx,
		tgClient: tgClient,
		storage:  storage,
		// search:          search,
		lastUserCommand: map[int64]string{},
	}
}

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	UserID          int64
	UserName        string
	UserLastName    string
	UserFirtsName   string
	UserDisplayName string
	IsBot           bool
	IsCallback      bool
	CallbackMsgID   string
}

func (s *Model) GetCtx() context.Context {
	return s.ctx
}

func (s *Model) SetCtx(ctx context.Context) {
	s.ctx = ctx
}

// IncomingMessage Обработка входящего сообщения.
func (s *Model) IncomingMessage(msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(s.ctx, "IncomingMessage")
	s.ctx = ctx
	defer span.Finish()

	// lastUserCommand := s.lastUserCommand[msg.UserID]

	// Обнуление выбранной команды.
	s.lastUserCommand[msg.UserID] = ""

	// Распознавание стандартных команд.
	if isNeedReturn, err := checkBotCommands(s, msg); err != nil || isNeedReturn {
		return err
	}

	// Отправка ответа по умолчанию.
	return s.tgClient.SendMessage(txtUnknownCommand, msg.UserID)
}

// Распознавание стандартных команд бота.
func checkBotCommands(s *Model, msg Message) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(s.ctx, "checkBotCommands")
	s.ctx = ctx
	defer span.Finish()

	switch msg.Text {
	case "/start":
		displayName := msg.UserDisplayName
		if len(displayName) == 0 {
			displayName = msg.UserName
		}
		// Отображение команд стартовых действий.
		return true, s.tgClient.ShowInlineButtons(
			fmt.Sprintf(txtStart, displayName),
			btnStart,
			msg.UserID,
		)
	case "/help":
		return true, s.tgClient.SendMessage(txtHelp, msg.UserID)
	case "/search_personal_card":
		s.lastUserCommand[msg.UserID] = "/search_personal_card"
		return true, s.tgClient.SendMessage(
			txtCriterionChoose,
			msg.UserID,
		)
	case "/all_personal_cards":
		cards, _ := s.storage.ShowAllPersonalCards(ctx)
		return true, s.tgClient.SendMessage(fmt.Sprintf("%v", cards), msg.UserID)
	}
	// Команда не распознана.
	return false, nil
}
