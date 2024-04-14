package dialog

import (
	"context"
	"fmt"
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
	SendCards(cards []string, chatID int64) error
	SendDBDump() error
	StartDBJob(ctx context.Context)
	SendFile(chatID int64, file *tgbotapi.FileReader, currentTime string) error
	EditTextAndMarkup(
		msg Message,
		newText string,
		newMarkup *tgbotapi.InlineKeyboardMarkup,
	) error
	EditMarkup(msg Message, markup *tgbotapi.InlineKeyboardMarkup) error
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
	Text            string
	Data            string
	MsgID           int
	ChatID          int64
	UserID          int64
	BotName         string
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
	if isNeedReturn, err := CheckBotCommands(ctx, m, msg); err != nil || isNeedReturn {
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
		m.search.AddSearchData(msg.Text)
		cards, err := m.search.ProcessCards(ctx, m.storage)
		if err != nil {
			logger.Error("Ошибка в поиске карты", "ERROR", err)
		}
		m.search.Disable()
		return m.tgClient.SendCards(cards, msg.ChatID)
	}

	// Отправка ответа по умолчанию.
	return m.tgClient.SendMessage(txtUnknownCommand, msg.ChatID)
}

// CheckBotCommands распознавание стандартных команд бота.
func CheckBotCommands(ctx context.Context, m *Model, msg Message) (bool, error) {
	switch msg.Text {
	case "/start", fmt.Sprintf("/start@" + msg.BotName):
		if m.search.IsEnabled() {
			m.search.Disable()
		}
		return true, m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, msg.FirstName),
			msg.ChatID,
			&MarkupMainMenu,
		)
	case "/allpersonalcards":
		rawCards, err := m.storage.ShowAllPersonalCards(ctx)
		if err != nil {
			logger.Error("Ошибка в сборе всех персональных карточек", "ERROR", err)
		}
		cards := card.FormatCards(rawCards)
		return true, m.tgClient.SendCards(cards, msg.ChatID)
	case "/dump":
		return true, m.tgClient.SendDBDump()
		// case "/dice":
		// 	return true, tgbotapi.Dice.Emoji
	}
	return false, nil
}

func (m *Model) HandleButton(msg Message) error {
	button := msg.CallbackQuery.Data
	firstName := msg.CallbackQuery.From.FirstName
	previousMarkup := msg.CallbackQuery.Message.ReplyMarkup

	switch button {
	case BtnBack:
		m.search.Disable()
		ResetCriteriaButtons()
		return m.tgClient.EditTextAndMarkup(
			msg,
			fmt.Sprintf(txtMainMenu, firstName),
			&MarkupMainMenu,
		)
	case BtnSearchPerson:
		m.search.SetSearchScreen("personal_cards")
		return m.tgClient.EditTextAndMarkup(
			msg,
			txtCriterionChoose,
			&MarkupSearchPersonMenu,
		)
	case BtnSearchOrganization:
		m.search.SetSearchScreen("organization_cards")
		return m.tgClient.EditTextAndMarkup(
			msg,
			txtCriterionChoose,
			&MarkupSearchOrganizationMenu,
		)
	case BtnApply:
		lenCriterions := len(m.search.GetCriterions())
		if lenCriterions == 0 {
			return m.tgClient.EditTextAndMarkup(
				msg,
				txtNoCriteria,
				previousMarkup,
			)
			// TEST
		} else if lenCriterions == 1 {
			m.search.Enable()
			var alias string
			for key := range m.search.GetCriterions() {
				alias = key
			}
			return m.tgClient.EditTextAndMarkup(
				msg,
				fmt.Sprintf(txtCriteriaInput, alias),
				&MarkupCancelMenu,
			)
		}
	case BtnCancelSearch:
		m.search.Disable()
		ResetCriteriaButtons()
		return m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, firstName),
			msg.ChatID,
			&MarkupMainMenu,
		)
	case BtnMenu:
		return m.tgClient.SendMessageWithMarkup(
			fmt.Sprintf(txtMainMenu, firstName),
			msg.ChatID,
			&MarkupMainMenu,
		)
	case HandleCriterionButton(button, m.search):
		searchScreen := m.search.GetSearchScreen()
		markup := CreateSearchMenuMarkup(searchScreen)
		return m.tgClient.EditMarkup(
			msg,
			&markup,
		)
	}

	return nil
}
