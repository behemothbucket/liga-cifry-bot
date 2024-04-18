package dialog

import (
	"context"
	"fmt"
	"os"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/card"
	"telegram-bot/internal/model/db"
	"telegram-bot/internal/model/search"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Область "Константы и переменные": начало.

var (
	txtMainMenu       = "Привет, %v.\nМогу помочь найти карточку компетенций или организации."
	txtUnknownCommand = "К сожалению, данная команда мне неизвестна.\nДля начала работы введите\n/start"
	txtCardNotFound   = "Ничего не найдено... 🤷‍♂️"
	// txtReportWait      = "Ищу 🔎\nПожалуйста, подождите..."
	txtCriterionChoose = "Выберите критерии поиска для поиска, а затем нажмите *Применить* ✅."
	txtNoCriteria      = "❗️Не выбрано ни одного критерия поиска. Сначала выберите хотя-бы один критерий."
	txtCriteriaInput   = "Пожалуйста, введите *%v*."
)

// Область "Константы и переменные": конец.

// Область "Внешний интерфейс": начало.

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(msg Message, typeMessage string) error
	SendMessageWithMarkup(msg Message) error
	SendFile(msg Message) error
	SendMedia(msg Message) error
	SendMediaGroup(msg Message) error
	EditTextAndMarkup(msg Message) error
	EditMarkup(msg Message) error
	DeferMessage(msg Message)
	StartDBJob(ctx context.Context)
	SendDBDump() error
}

// Model Модель бота (клиент, хранилище, поиск)
type Model struct {
	ctx      context.Context
	tgClient MessageSender      // Клиент.
	storage  db.UserDataStorage // Хранилище пользовательской информации.
	search   search.Engine      // Поиск.
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей и параметров поиска.
func New(
	ctx context.Context,
	tgClient MessageSender,
	storage db.UserDataStorage,
	searchEngine search.Engine,
) *Model {
	return &Model{
		ctx:      ctx,
		tgClient: tgClient,
		storage:  storage,
		search:   searchEngine,
	}
}

// Message Структура сообщения для обработки.
type Message struct {
	IsCommand      bool
	MsgID          int
	ChatID         int64
	Text           string
	NewText        string
	BotName        string
	FirstName      string
	UserName       string
	Caption        string
	Type           string
	FilePaths      []string
	Cards          []string
	NewChatMembers []tgbotapi.User
	CallbackQuery  *tgbotapi.CallbackQuery
	LeftChatMember *tgbotapi.User
	Markup         tgbotapi.InlineKeyboardMarkup
	File           *tgbotapi.FileReader
	Photo          tgbotapi.RequestFileData
}

func (m *Model) GetCtx() context.Context {
	return m.ctx
}

func (m *Model) SetCtx(ctx context.Context) {
	m.ctx = ctx
}

// HandleMessage Обработка входящего сообщения.
func (m *Model) HandleMessage(msg Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch {
	case msg.IsCommand:
		HandleBotCommands(ctx, m, msg)
	case len(msg.NewChatMembers) != 0:
		if err := m.storage.JoinGroup(ctx, &msg.NewChatMembers[0]); err != nil {
			logger.Error("Ошибка в добавлении нового пользователя", "ERROR", err)
		}
	case msg.LeftChatMember != nil:
		if err := m.storage.LeaveGroup(ctx, msg.LeftChatMember); err != nil {
			logger.Error("Ошибка в исключении нового пользователя", "ERROR", err)
		}
	case m.search.IsEnabled():
		m.search.AddSearchData(msg.Text)
		cards, err := m.search.ProcessCards(ctx, m.storage)
		if cards == nil {
			logger.Info(
				"Не найдено ни одной записи по данному запросу",
			)
			msg.Text = txtCardNotFound
			msg.Markup = MarkupCancelMenu
			msg.Type = "SendMessageWithMarkup"
			m.tgClient.DeferMessage(msg)
		}
		if err != nil {
			logger.Error("Ошибка в поиске карты", "ERROR", err)
		}
		logger.Debug(fmt.Sprintf("%v", cards))
		msg.Cards = cards
		msg.Markup = MarkupCardMenu
		msg.Type = "SendMessageWithMarkup"
		m.search.Disable()
		m.tgClient.DeferMessage(msg)
	default:
		msg.Text = txtUnknownCommand
		msg.Type = ""
		m.tgClient.DeferMessage(msg)
	}
}

// CheckBotCommands распознавание стандартных команд бота.
func HandleBotCommands(ctx context.Context, m *Model, msg Message) {
	// TEST
	testChatID := int64(5587823077)
	// testChatID := int64(155401792)
	switch msg.Text {
	case "/start", fmt.Sprintf("/start@" + msg.BotName):
		if m.search.IsEnabled() {
			m.search.Disable()
		}
		msg.Text = fmt.Sprintf(txtMainMenu, msg.FirstName)
		msg.Markup = MarkupMainMenu
		msg.Type = "SendMessageWithMarkup"
		m.tgClient.DeferMessage(msg)
	case "/allpersonalcards":
		rawCards, err := m.storage.ShowAllPersonalCards(ctx)
		if err != nil {
			logger.Error("Ошибка в сборе всех персональных карточек", "ERROR", err)
		}
		cards := card.FormatCards(rawCards)
		msg.Cards = cards
		msg.Type = "SendMessageWithMarkup"
		m.tgClient.DeferMessage(msg)
	case "/dump":
		if err := m.tgClient.SendDBDump(); err != nil {
			logger.Error("Ошибка при выгрузка дампа БД", "ERROR", err)
		}
	case "/cat":
		file, _ := os.Open("./img/cat.jpg")
		reader := tgbotapi.FileReader{Name: file.Name(), Reader: file}
		msg.ChatID = testChatID
		msg.File = &reader
		msg.Caption = "Здарова ептить"
		msg.Type = "SendMedia"
		m.tgClient.DeferMessage(msg)
	case "/cats":
		paths := []string{"./img/cat.jpg", "./img/cat.jpg", "./img/cat.jpg"}
		msg.FilePaths = paths
		msg.ChatID = testChatID
		msg.Caption = "Бэйби"
		m.tgClient.DeferMessage(msg)
	}
}

func (m *Model) HandleButton(msg Message) {
	button := msg.CallbackQuery.Data
	firstName := msg.CallbackQuery.From.FirstName
	previousMarkup := *msg.CallbackQuery.Message.ReplyMarkup

	switch button {
	case BtnBack:
		m.search.Disable()
		ResetCriteriaButtons()
		msg.Type = "EditTextAndMarkup"
		msg.NewText = fmt.Sprintf(txtMainMenu, firstName)
		msg.Markup = MarkupMainMenu
	case BtnSearchPerson:
		m.search.SetSearchScreen("personal_cards")
		msg.Type = "EditTextAndMarkup"
		msg.NewText = txtCriterionChoose
		msg.Markup = MarkupSearchPersonMenu
	case BtnSearchOrganization:
		m.search.SetSearchScreen("organization_cards")
		m.search.SetSearchScreen("personal_cards")
		msg.Type = "EditTextAndMarkup"
		msg.NewText = txtCriterionChoose
		msg.Markup = MarkupSearchOrganizationMenu
	case BtnApply:
		lenCriterions := len(m.search.GetCriterions())
		if lenCriterions == 0 {
			msg.Type = "EditTextAndMarkup"
			msg.NewText = txtNoCriteria
			msg.Markup = previousMarkup
		} else if lenCriterions == 1 {
			m.search.Enable()
			var alias string
			for key := range m.search.GetCriterions() {
				alias = key
			}
			msg.Type = "SendMessageWithMarkup"
			msg.Text = fmt.Sprintf(txtCriteriaInput, alias)
			msg.Markup = MarkupCancelMenu
		}
	case BtnCancelSearch:
		m.search.Disable()
		ResetCriteriaButtons()
		msg.Type = "SendMessageWithMarkup"
		msg.Text = fmt.Sprintf(txtMainMenu, firstName)
		msg.Markup = MarkupMainMenu
	case BtnMenu:
		m.search.Disable()
		ResetCriteriaButtons()
		msg.Type = "SendMessageWithMarkup"
		msg.Text = fmt.Sprintf(txtMainMenu, firstName)
		msg.Markup = MarkupMainMenu
	case HandleCriterionButton(button, m.search):
		searchScreen := m.search.GetSearchScreen()
		msg.Type = "EditMarkup"
		msg.Markup = CreateSearchMenuMarkup(searchScreen)
	}
	m.tgClient.DeferMessage(msg)
}
