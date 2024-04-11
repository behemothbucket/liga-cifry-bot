package dialog

import (
	"context"
	"fmt"
	"telegram-bot/internal/helpers/markup"
	"telegram-bot/internal/logger"

	db "telegram-bot/internal/model/db"
	search "telegram-bot/internal/model/search"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
)

// Область "Константы и переменные": начало.

var (
	txtMainMenu = markup.EscapeForMarkdown(
		"Привет, %v.\nМогу помочь найти карточку компетенций или организации.",
	)
	txtUnknownCommand = markup.EscapeForMarkdown(
		"К сожалению, данная команда мне неизвестна.\nДля начала работы введите\n/start",
	)
	// txtReportError     = "Не удалось получить данные."
	// txtReportWait      = "Ищу 🔎\nПожалуйста, подождите..."
	txtCriterionChoose = markup.EscapeForMarkdown(
		"Выберите критерии поиска для поиска, а затем нажмите *Применить* ✅.",
	)
	txtNoCriteria = markup.EscapeForMarkdown(
		"❗️Не выбрано ни одного критерия поиска. Сначала выберите хотя-бы один критерий.",
	)
	txtCriteriaInput = markup.EscapeForMarkdown(
		"Пожалуйста, введите *%v*.",
	)
)

// Область "Константы и переменные": конец.

// Область "Внешний интерфейс": начало.

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, chatID int64) error
	SendMessageWithMarkup(text string, chatID int64, markup tgbotapi.InlineKeyboardMarkup) error
	ShowInlineButtons(
		chatID int64,
		msgID int,
		text string,
		markup tgbotapi.InlineKeyboardMarkup,
	) error
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
	MsgID           int
	ChatID          int64
	UserID          int64
	FirstName       string
	CallbackQuery   *tgbotapi.CallbackQuery
	NewChatMembers  []tgbotapi.User
	LeftChatMembers *tgbotapi.User
}

func (m *Model) GetCtx() context.Context {
	return m.ctx
}

func (m *Model) SetCtx(ctx context.Context) {
	m.ctx = ctx
}

// HandleMessage Обработка входящего сообщения.
func (m *Model) HandleMessage(msg Message) error {
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

	// Режим поиска
	if m.search.GetMode() != "" {
		card, _ := m.storage.FindCard(ctx, m.search.GetCriterions()[0], msg.Text)
		logger.Debug(card)
		// return m.tgClient.SendMessage(
		// 	card,
		// 	msg.ChatID,
		// )
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
		m.search.Disable()
		displayName := msg.FirstName
		if len(displayName) == 0 {
			displayName = msg.FirstName
		}
		// Отображение команд стартовых действий.
		return true, m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, displayName),
			msg.ChatID,
			markupMainMenu,
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
	mode := m.search.GetMode()

	switch button {
	case btnBack:
		m.search.Disable()
		displayName := msg.CallbackQuery.From.FirstName
		return m.tgClient.ShowInlineButtons(
			msg.ChatID,
			msg.MsgID,
			fmt.Sprintf(txtMainMenu, displayName),
			markupMainMenu,
		)
	case btnSearchPerson:
		m.search.SetMode("person")
		return m.tgClient.ShowInlineButtons(
			msg.ChatID,
			msg.MsgID,
			txtCriterionChoose,
			markupSearchPersonMenu,
		)
	case btnSearchOrganization:
		m.search.SetMode("organization")
		return m.tgClient.ShowInlineButtons(
			msg.ChatID,
			msg.MsgID,
			txtCriterionChoose,
			markupSearchOrganizationMenu,
		)
	case btnApply:
		lenCriterions := len(m.search.GetCriterions())
		if lenCriterions == 0 {
			markup := CreateSearchMenuMarkup(mode)
			return m.tgClient.ShowInlineButtons(
				msg.ChatID,
				msg.MsgID,
				txtNoCriteria,
				markup,
			)
		} else if lenCriterions == 1 {
			markup := CreateCancelMenuMarkup()
			return m.tgClient.ShowInlineButtons(
				msg.ChatID,
				msg.MsgID,
				fmt.Sprintf(txtCriteriaInput, m.search.GetCriterions()[0]),
				markup,
			)
		}
	case IsCriterionButton(button, mode):
		toggleCriterionButton(button, m.search)
		logger.Debug(fmt.Sprintf("%v", m.search.GetCriterions()))
		// NOTE можно ли по ссылке менять?
		markup := CreateSearchMenuMarkup(mode)
		return m.tgClient.ShowInlineButtons(
			msg.ChatID,
			msg.MsgID,
			txtCriterionChoose,
			markup,
		)
	}

	return nil
}
