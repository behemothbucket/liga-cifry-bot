package messages

import (
	"context"
	"fmt"
	types "telegram-bot/internal/model/bottypes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
)

// Область "Константы и переменные": начало.

const (
	txtStart           = "Привет, %v.\nЯ - бот, помогающий искать карточки компетенций и организаций для проекта <b>Лига Цифры</b>."
	txtUnknownCommand  = "К сожалению, данная команда мне неизвестна.\nДля начала работы введите\n/start"
	txtReportError     = "Не удалось получить данные."
	txtReportWait      = "Ищу 🔎\nПожалуйста, подождите..."
	txtCriterionChoose = "Выберите критерии поиска для поиска , а затем нажмите <b>Применить</b>."
	txtCatEmpty        = "Не выбрано ни одного критерия поиска. Сначала выберите хотя-бы один критерий."
)

// Кнопки.
var (
	searchButtons = map[string][]string{
		"user": {
			"ФИО",
			"Город",
			"Организация",
			"Должность",
			"Компетенции",
			"Направления сотрудничества",
		},
		"university": {
			"Организация",
			"Структурное подразделение",
			"Город",
			"Направления сотрудничества",
			"«Приоритет-2030»",
			"Членство в консорциуме",
			"Разработки отвечественного ПО",
			"Лабораторные площадки и НОЦ",
			"Компетенции",
		},
	}
	searchPersonButton         = "🔍 Индивидульные карточки"
	searchOrganizationButton   = "🔍 Карточки организаций"
	backButton                 = "⬅️ Назад"
	menuButton                 = "↩️ Меню"
	cancelSearchButton         = "❌ Отменить поиск"
	applyButton                = "🆗 Применить"
	searchButton               = "🔍 Искать"
	printFirstPersonalCard     = "⚠️Персональная карточка⚠️"
	printAllPersonalCards      = "⚠️Все персональные карточки⚠️"
	printFirstOrganizationCard = "⚠️Карточка организации⚠️"
	loadMoreButton             = "⏬ Загрузить еще 5"

	toggleButtonPrefix = "✅ "
)

// Кнопки поиска в главном меню.
var btnSearch = []types.TgRowButtons{
	{
		types.TgInlineButton{
			DisplayName: searchPersonButton,
		},
	},
	{
		types.TgInlineButton{
			DisplayName: searchOrganizationButton,
		},
	},
}

var btnSearchPerson = []types.TgRowButtons{
	{
		types.TgInlineButton{
			DisplayName: "ФИО",
		},
		types.TgInlineButton{
			DisplayName: "Город",
		},
		types.TgInlineButton{
			DisplayName: "Организация",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: "Должность",
		},
		types.TgInlineButton{
			DisplayName: "Компетенции",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: "Направления сотрудничества",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: backButton,
		},
		types.TgInlineButton{
			DisplayName: applyButton,
		},
	},
}

var btnSearchOrganization = []types.TgRowButtons{
	{
		types.TgInlineButton{
			DisplayName: "Структурное подразделение",
		},
		types.TgInlineButton{
			DisplayName: "«Приоритет-2030»",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: "Членство в консорциуме",
		},
		types.TgInlineButton{
			DisplayName: "Компетенции",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: "Лабораторные площадки и НОЦ",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: "Разработки отвечественного ПО",
		},
	},
	{
		types.TgInlineButton{
			DisplayName: backButton,
		},
		types.TgInlineButton{
			DisplayName: applyButton,
		},
	},
}

var (
	mainMenuMarkup               = getMainMenuMarkup()
	searchPersonMenuMarkup       = getSearchPersonCardMenu()
	searchOrganizationMenuMarkup = getSearchOrganizationCardMenu()
)

func createNumericKeyboard(buttons []types.TgRowButtons) tgbotapi.InlineKeyboardMarkup {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for i := 0; i < len(buttons); i++ {
		tgRowButtons := buttons[i]
		keyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(tgRowButtons))
		for j := 0; j < len(tgRowButtons); j++ {
			tgInlineButton := tgRowButtons[j]
			keyboard[i][j] = tgbotapi.NewInlineKeyboardButtonData(
				tgInlineButton.DisplayName,
				tgInlineButton.DisplayName,
			)
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func getMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return createNumericKeyboard(btnSearch)
}

func getSearchPersonCardMenu() tgbotapi.InlineKeyboardMarkup {
	return createNumericKeyboard(btnSearchPerson)
}

func getSearchOrganizationCardMenu() tgbotapi.InlineKeyboardMarkup {
	return createNumericKeyboard(btnSearchOrganization)
}

// Область "Константы и переменные": конец.

// Область "Внешний интерфейс": начало.

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, userID int64) error
	ShowInlineButtons(text string, markup tgbotapi.InlineKeyboardMarkup, userID int64) error
}

// UserDataStorage Интерфейс для работы с хранилищем данных.
type UserDataStorage interface {
	JoinGroup(ctx context.Context, u *tgbotapi.User) error
	LeaveGroup(ctx context.Context, u *tgbotapi.User) error
	CheckIfUserExist(ctx context.Context, userID int64) (bool, error)
	ShowAllPersonalCards(ctx context.Context) (pc []PersonalCard, err error)
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
	}
}

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	ChatID          int64
	UserID          int64
	FirstName       string
	NewChatMembers  []tgbotapi.User
	LeftChatMembers *tgbotapi.User
	CallbackQuery   *tgbotapi.CallbackQuery
	IsCallback      bool
	CallbackMsgID   string
}

func (m *Model) GetCtx() context.Context {
	return m.ctx
}

func (m *Model) SetCtx(ctx context.Context) {
	m.ctx = ctx
}

// IncomingMessage Обработка входящего сообщения.
func (m *Model) IncomingMessage(msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(m.ctx, "IncomingMessage")
	m.ctx = ctx
	defer span.Finish()

	// Распознавание стандартных команд.
	if isNeedReturn, err := checkBotCommands(m, msg); err != nil || isNeedReturn {
		return err
	}

	// Новый участник группы
	if len(msg.NewChatMembers) != 0 {
		return m.storage.JoinGroup(ctx, &msg.NewChatMembers[0])
	}

	// Участник вышел или был удален
	if msg.LeftChatMembers != nil {
		return m.storage.LeaveGroup(ctx, msg.LeftChatMembers)
	}

	// Отправка ответа по умолчанию.
	return m.tgClient.SendMessage(txtUnknownCommand, msg.ChatID)
}

// Распознавание стандартных команд бота.
func checkBotCommands(m *Model, msg Message) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(m.ctx, "checkBotCommands")
	m.ctx = ctx
	defer span.Finish()

	switch msg.Text {
	case "/start":
		displayName := msg.FirstName
		if len(displayName) == 0 {
			displayName = msg.FirstName
		}
		// Отображение команд стартовых действий.
		return true, m.tgClient.ShowInlineButtons(
			fmt.Sprintf(txtStart, displayName),
			mainMenuMarkup,
			msg.UserID,
		)
	}
	// Команда не распознана.
	return false, nil
}

func (m *Model) HandleButton(msg Message) error {
	span, ctx := opentracing.StartSpanFromContext(m.ctx, "HandleButton")
	m.ctx = ctx
	defer span.Finish()

	button := msg.CallbackQuery.Data

	switch button {
	case backButton:
		displayName := msg.CallbackQuery.From.FirstName
		return m.tgClient.ShowInlineButtons(
			fmt.Sprintf(txtStart, displayName),
			mainMenuMarkup,
			msg.CallbackQuery.Message.Chat.ID,
		)
	case searchPersonButton:
		return m.tgClient.ShowInlineButtons(
			txtCriterionChoose,
			searchPersonMenuMarkup,
			msg.CallbackQuery.Message.Chat.ID,
		)
	case searchOrganizationButton:
		return m.tgClient.ShowInlineButtons(
			txtCriterionChoose,
			searchOrganizationMenuMarkup,
			msg.CallbackQuery.Message.Chat.ID,
		)
	}

	return nil
}
