package dialog

import (
	"context"
	"fmt"
	"telegram-bot/internal/model/card/person"
	"time"

	db "telegram-bot/internal/model/db"
	search "telegram-bot/internal/model/search"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
)

// Область "Константы и переменные": начало.

var (
	txtMainMenu       = "Привет, %v.\nМогу помочь найти карточку компетенций или организации."
	txtUnknownCommand = "К сожалению, данная команда мне неизвестна.\nДля начала работы введите\n/start"
	// txtReportError     = "Не удалось получить данные."
	// txtReportWait      = "Ищу 🔎\nПожалуйста, подождите..."
	txtCriterionChoose = "Выберите критерии поиска для поиска, а затем нажмите *Применить* ✅."
	txtNoCriteria      = "❗️Не выбрано ни одного критерия поиска. Сначала выберите хотя-бы один критерий."
	txtCriteriaInput   = "Пожалуйста, введите *%v*."
)

// Область "Константы и переменные": конец.

// Область "Внешний интерфейс": начало.

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, chatID int64) error
	SendMessageWithMarkup(text string, chatID int64, markup *tgbotapi.InlineKeyboardMarkup) error
	EditTextAndMarkup(
		msg Message,
		newText string,
		newMarkup *tgbotapi.InlineKeyboardMarkup,
	) error
	EditMarkup(msg Message, markup *tgbotapi.InlineKeyboardMarkup) error
	// EditText(chatID int64, msgID int, text string) error
}

// Model Модель бота (клиент, хранилище, поиск)
type Model struct {
	ctx      context.Context
	tgClient MessageSender       // Клиент.
	storage  db.UserDataStorage  // Хранилище пользовательской информации.
	search   search.SearchEngine // Поиск.
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей и параметров поиска.
func New(
	ctx context.Context,
	tgClient MessageSender,
	storage db.UserDataStorage,
	searchEngine search.SearchEngine,
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
	Text            string
	Data            string
	MsgID           int
	ChatID          int64
	UserID          int64
	FirstName       string
	CallbackQuery   *tgbotapi.CallbackQuery
	NewChatMembers  []tgbotapi.User
	LeftChatMembers *tgbotapi.User
	Markup          *tgbotapi.InlineKeyboardMarkup
}

func (m *Model) GetCtx() context.Context {
	return m.ctx
}

func (m *Model) SetCtx(ctx context.Context) {
	m.ctx = ctx
}

// HandleMessage Обработка входящего сообщения.
func (m *Model) HandleMessage(msg Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Распознавание стандартных команд.
	if isNeedReturn, err := CheckBotCommands(m, msg); err != nil || isNeedReturn {
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

	// Режим поиска
	if m.search.IsEnabled() {
		return m.ProcessSearch(ctx, msg)
	}

	// Отправка ответа по умолчанию.
	return m.tgClient.SendMessage(txtUnknownCommand, msg.ChatID)
}

// Распознавание стандартных команд бота.
func CheckBotCommands(m *Model, msg Message) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(m.ctx, "checkBotCommands")
	m.ctx = ctx
	defer span.Finish()

	switch msg.Text {
	case "/start":
		m.search.Disable()
		// Отображение команд стартовых действий.
		return true, m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, msg.FirstName),
			msg.ChatID,
			&markupMainMenu,
		)
	}
	// Команда не распознана.
	return false, nil
}

func (m *Model) ProcessSearch(ctx context.Context, msg Message) error {
	rawCard, err := m.storage.FindCard(ctx, m.search.GetCriterions()[0], msg.Text)
	if err != nil {
		return err
	}

	card := person.MarkupCard(&rawCard)

	return m.tgClient.SendMessage(
		card,
		msg.ChatID,
	)
}

func (m *Model) HandleButton(msg Message) error {
	button := msg.CallbackQuery.Data
	searchScreen := m.search.GetSearchScreen()
	firstName := msg.CallbackQuery.From.FirstName
	previousMarkup := msg.CallbackQuery.Message.ReplyMarkup

	switch button {
	case btnBack:
		m.search.Disable()
		ResetCriteriaButtons()
		return m.tgClient.EditTextAndMarkup(
			msg,
			fmt.Sprintf(txtMainMenu, firstName),
			&markupMainMenu,
		)
	case btnSearchPerson:
		m.search.SetSearchScreen("person")
		return m.tgClient.EditTextAndMarkup(
			msg,
			txtCriterionChoose,
			&markupSearchPersonMenu,
		)
	case btnSearchOrganization:
		m.search.SetSearchScreen("organization")
		return m.tgClient.EditTextAndMarkup(
			msg,
			txtCriterionChoose,
			&markupSearchOrganizationMenu,
		)
	case btnApply:
		m.search.Enable()
		lenCriterions := len(m.search.GetCriterions())
		if lenCriterions == 0 {
			return m.tgClient.EditTextAndMarkup(
				msg,
				txtNoCriteria,
				previousMarkup,
			)
			// TEST
		} else if lenCriterions == 1 {
			markup := markupCancelMenu
			return m.tgClient.EditTextAndMarkup(
				msg,
				fmt.Sprintf(txtCriteriaInput, m.search.GetCriterions()[0]),
				&markup,
			)
		}
	case btnCancelSearch:
		m.search.Disable()
		m.search.ResetSearchCriterions()
		ResetCriteriaButtons()
		return m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, firstName),
			msg.ChatID,
			&markupMainMenu,
		)
	case HandleCriterionButton(button, m.search):
		markup := CreateSearchMenuMarkup(searchScreen)
		return m.tgClient.EditMarkup(
			msg,
			&markup,
		)
	}

	return nil
}
