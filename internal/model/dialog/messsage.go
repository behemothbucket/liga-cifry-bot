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
	SendMessage(msg Message) error
	// SendMessageWithMarkup(
	// 	chatID int64,
	// 	text string,
	// 	markup *tgbotapi.InlineKeyboardMarkup,
	// ) DeferredMessage
	// SendCards(chatID int64, cards []string) DeferredMessage
	// SendFile(chatID int64, file *tgbotapi.FileReader, currentTime string) DeferredMessage
	// SendMedia(chatID int64, file *tgbotapi.FileReader, caption string) DeferredMessage
	// SendMediaGroup(chatID int64, paths []string, caption string) DeferredMessage
	EditTextAndMarkup(
		msg Message,
		newText string,
		newMarkup *tgbotapi.InlineKeyboardMarkup,
	) error
	EditMarkup(msg Message, markup *tgbotapi.InlineKeyboardMarkup) error
	DeferUpdate(update tgbotapi.Update)
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
	Text           string
	ChatID         int64
	IsCommand      bool
	BotName        string
	FirstName      string
	MsgID          int
	CallbackQuery  *tgbotapi.CallbackQuery
	NewChatMembers []tgbotapi.User
	LeftChatMember *tgbotapi.User
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

	switch {
	case msg.IsCommand:
		return HandleBotCommands(ctx, m, msg)
	case len(msg.NewChatMembers) != 0:
		return m.storage.JoinGroup(ctx, &msg.NewChatMembers[0])
	case msg.LeftChatMembers != nil:
		return m.storage.LeaveGroup(ctx, msg.LeftChatMembers)
	case m.search.IsEnabled():
		m.search.AddSearchData(msg.Text)
		cards, err := m.search.ProcessCards(ctx, m.storage)
		if cards == nil {
			logger.Info(
				"Не найдено ни одной записи по данному запросу",
			)
			return m.tgClient.SendMessageWithMarkup(
				msg.ChatID,
				txtCardNotFound,
				&MarkupCardMenu,
			)
		}
		if err != nil {
			logger.Error("Ошибка в поиске карты", "ERROR", err)
		}
		m.search.Disable()
		return m.tgClient.SendCards(msg.ChatID, cards)
	default:
		return m.tgClient.SendMessage(msg.ChatID, txtUnknownCommand)
	}
}

// CheckBotCommands распознавание стандартных команд бота.
func HandleBotCommands(ctx context.Context, m *Model, msg Message) error {
	// TEST
	testChatID := int64(5587823077)
	// testChatID := int64(155401792)
	switch msg.Text {
	case "/start", fmt.Sprintf("/start@" + msg.BotName):
		if m.search.IsEnabled() {
			m.search.Disable()
		}
		return m.tgClient.SendMessageWithMarkup(
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, msg.FirstName),
			&MarkupMainMenu,
		)
	case "/allpersonalcards":
		rawCards, err := m.storage.ShowAllPersonalCards(ctx)
		if err != nil {
			logger.Error("Ошибка в сборе всех персональных карточек", "ERROR", err)
		}
		cards := card.FormatCards(rawCards)
		return m.tgClient.SendCards(msg.ChatID, cards)
	case "/dump":
		return m.tgClient.SendDBDump()
	case "/cat":
		file, _ := os.Open("./img/cat.jpg")
		reader := tgbotapi.FileReader{Name: file.Name(), Reader: file}
		return m.tgClient.SendMedia(testChatID, &reader, "Здарова ептить")
	case "/cats":
		paths := []string{"./img/cat.jpg", "./img/cat.jpg", "./img/cat.jpg"}
		return m.tgClient.SendMediaGroup(testChatID, paths, "Бэйби")
	}

	return nil
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
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, firstName),
			&MarkupMainMenu,
		)
	case BtnMenu:
		m.search.Disable()
		ResetCriteriaButtons()
		return m.tgClient.SendMessageWithMarkup(
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, firstName),
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
