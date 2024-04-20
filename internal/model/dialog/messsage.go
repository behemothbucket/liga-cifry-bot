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

var (
	txtMainMenu        = "👋 Привет, <b>%v</b>.\nМогу помочь найти карточку компетенций или организации."
	txtUnknownMessage  = "💬 <b>К сожалению, данная команда мне неизвестна.</b>\n\nДля начала работы введите\n<u>/start</u>."
	txtCardNotFound    = "Ничего не найдено... 🤷‍♂️"
	txtCriterionChoose = "💬 <b>Выберите критерии поиска и нажмите Применить ✅</b>"
	txtNoCriteria      = "❗️<b>Не выбрано ни одного критерия поиска</b>"
	txtCriteriaInput   = "Пожалуйста, введите <b>%v</b>"
	themeLink          = "<a href='https://t.me/addtheme/Liga_Cifry'>🎨 Офицальная цветовая тема Лиги Цифры!</a>"
)

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithMarkup(chatID int64, text string, markup *tgbotapi.ReplyKeyboardMarkup) error
	SendKeyboard(chatID int64, text string, markup *tgbotapi.ReplyKeyboardMarkup) error
	SendCards(chatID int64, cards []string) error
	SendDBDump() error
	StartDBJob(ctx context.Context)
	SendFile(chatID int64, file *tgbotapi.FileReader, currentTime string) error
	SendMedia(chatID int64, file *tgbotapi.FileReader, caption string) error
	SendMediaGroup(chatID int64, paths []string, caption string) error
	DeferMessageWithMarkup(msg Message)
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
	Data           string
	MsgID          int
	Markup         *tgbotapi.ReplyKeyboardMarkup
	ChatID         int64
	UserID         int64
	BotName        string
	FirstName      string
	IsCommand      bool
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

	if msg.IsCommand {
		return handleCommand(ctx, m, msg)
	}

	if len(msg.NewChatMembers) != 0 {
		return m.storage.JoinGroup(ctx, &msg.NewChatMembers[0])
	}

	if msg.LeftChatMember != nil {
		return m.storage.LeaveGroup(ctx, msg.LeftChatMember)
	}

	if m.search.IsEnabled() {
		m.search.AddSearchData(msg.Text)
		cards, err := m.search.ProcessCards(ctx, m.storage)
		if len(cards) == 0 {
			m.search.Disable()
			return m.tgClient.SendKeyboard(
				msg.ChatID,
				txtCardNotFound,
				&CardKeyboard,
			)
		}
		if err != nil {
			logger.Error("Ошибка в поиске карты", "ERROR", err)
		}
		m.search.Disable()
		return m.tgClient.SendCards(msg.ChatID, cards)
	}

	if isButton, err := handleButton(msg, m); err != nil || isButton {
		return err
	}

	return m.tgClient.SendMessage(msg.ChatID, txtUnknownMessage)
}

// CheckBotCommands распознавание стандартных команд бота.
func handleCommand(ctx context.Context, m *Model, msg Message) error {
	// TEST
	testChatID := int64(5587823077)
	// testChatID := int64(155401792)
	switch msg.Text {
	case "/start":
		if m.search.IsEnabled() {
			m.search.Disable()
		}
		ResetCriteriaButtons()
		return m.tgClient.SendKeyboard(
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, msg.FirstName),
			&MainKeyboard,
		)
	case "/theme":
		return m.tgClient.SendMessage(msg.ChatID, themeLink)

	case "/allpersonalcards":
		rawCards, err := m.storage.ShowAllPersonalCards(ctx)
		if err != nil {
			logger.Error("Ошибка в сборе всех персональных карточек", "ERROR", err)
		}
		cards := card.FormatPersonCards(rawCards)
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

func handleButton(msg Message, m *Model) (bool, error) {
	button := msg.Text
	firstName := msg.FirstName

	switch button {
	case BtnBack, BtnMenu:
		m.search.Disable()
		ResetCriteriaButtons()
		return true, m.tgClient.SendKeyboard(
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, msg.FirstName),
			&MainKeyboard,
		)
	case BtnSearchPerson:
		m.search.SetSearchScreen("personal_cards")
		return true, m.tgClient.SendKeyboard(msg.ChatID, txtCriterionChoose, &PersonKeyboard)
	case BtnSearchOrganization:
		m.search.SetSearchScreen("organization_cards")
		return true, m.tgClient.SendKeyboard(msg.ChatID, txtCriterionChoose, &OrganizationKeyboard)
	case BtnApply:
		lenCriterions := len(m.search.GetCriterions())
		if lenCriterions == 0 {
			return true, m.tgClient.SendMessage(msg.ChatID, txtNoCriteria)
			// TEST
		} else if lenCriterions == 1 {
			m.search.Enable()
			var alias string
			for key := range m.search.GetCriterions() {
				alias = key
			}
			return true, m.tgClient.SendKeyboard(
				msg.ChatID,
				fmt.Sprintf(txtCriteriaInput, alias),
				&CancelKeyboard,
			)
		}
	case BtnCancelSearch:
		m.search.Disable()
		ResetCriteriaButtons()
		return true, m.tgClient.SendMessageWithMarkup(
			msg.ChatID,
			fmt.Sprintf(txtMainMenu, firstName),
			&MainKeyboard,
		)
	case HandleCriterionButton(button, m.search):
		searchScreen := m.search.GetSearchScreen()
		markup := CreateSearchMenu(searchScreen)
		return true, m.tgClient.SendKeyboard(
			msg.ChatID,
			txtCriterionChoose,
			&markup,
		)
	}

	return false, nil
}
