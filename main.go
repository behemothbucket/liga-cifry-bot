package main

import (
	"bufio"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	mainMenuDescription   = "<b>Меню</b>\n\nТекст Текст Текст Текст"
	searchMenuDescription = "<b>Выберите критерии поиска:</b>"
)

var (
	searchButtons = map[string][]string{
		"user": {
			"ФИО",
			"Город",
			"ВУЗ",
			"Должность",
			"Компетенции",
			"Запросы на помощь/консультацию",
			"Сотрудничество",
		},
		"university": {
			"ВУЗ",
			"Приоритет 2030",
			"Кампус мирового уровня",
			"Наличие собственных разработок...",
			"Наличие лабораторных площадок...",
			"Компетенции",
			"Запросы на помощь/консультацию",
			"Сотрудничество",
		},
	}

	searchMode      = false
	searchCriterias = map[string]string{}
	currentCriteria = ""

	searchUserButton       = "Поиск карточки пользователя"
	searchUniversityButton = "Поиск карточки университета"
	backButton             = "⬅️ Назад"
	cancelButton           = "❌ Отмена"
	applyButton            = "📝 Применить"
	searchButton           = "🔍 Искать"
	miroButton             = "Miro"

	currentSearchScreen = ""

	bot *tgbotapi.BotAPI
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

func main() {
	telegramBotToken, exists := os.LookupEnv("TOKEN")

	if !exists {
		log.Print("Токен не обнаружен")
	}

	var err error

	bot, err = tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	log.Println("Сервер запущен. Нажмите Enter для остановки...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(update.Message)
		break
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	if user == nil {
		return
	}

	userName := user.UserName
	firstName := user.FirstName
	lastName := user.LastName
	userID := user.ID

	if lastName != "" {
		lastName = " " + lastName
	}

	log.Printf("https://t.me/%s [%d] (%s%s) написал(а) '%s'", userName, userID, firstName, lastName, text)

	var err error

	if reflect.TypeOf(message.Text).Kind() == reflect.String && message.Text != "" {
		if strings.HasPrefix(text, "/") {
			err = handleCommand(message.Chat.ID, text)
		} else if searchMode {
			if len(searchCriterias) == 0 {
				msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Вы заполнили все критерии</b>✅")
				msg.ParseMode = tgbotapi.ModeHTML
				_, err = bot.Send(msg)
				err = SendMenu(message.Chat.ID)
				searchMode = false
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "<b>Ответ принят</b>\nЯ пока что в разработке...")
			msg.ParseMode = tgbotapi.ModeHTML
			_, err = bot.Send(msg)
			delete(searchCriterias, currentCriteria)
		} else {
			err = SendMenu(message.Chat.ID)
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Файлы, фото/видео и другие медиа не принимаются ❌")
		_, err = bot.Send(msg)
		err = SendMenu(message.Chat.ID)
	}

	if err != nil {
		log.Printf("Ошибка: %s", err.Error())
	}
}

func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		err = SendMenu(chatId)
		break
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := getMainMenuMarkup()
	message := query.Message

	if query.Data == searchUserButton {
		text = searchMenuDescription
		markup = getUserSearchMenuMarkup()
		currentSearchScreen = "user"

	} else if query.Data == searchUniversityButton {
		text = searchMenuDescription
		markup = getUniversitySearchMenuMarkup()
		currentSearchScreen = "university"
	} else if query.Data == backButton {
		text = mainMenuDescription
		markup = getMainMenuMarkup()
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
	} else if query.Data == applyButton {
		text = getCriteria()
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(cancelButton, cancelButton)),
		)
		searchMode = true
	} else if query.Data == cancelButton {
		SendMenu(message.Chat.ID)
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
		searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		bot.Send(callbackCfg)
		return
	} else if criteriaButtonIsClicked(query.Data) {
		toggleButtonCheck(query.Data)
		text = searchMenuDescription
		if currentSearchScreen == "user" {
			markup = getUserSearchMenuMarkup()
		} else {
			markup = getUniversitySearchMenuMarkup()
		}
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func criteriaButtonIsClicked(button string) bool {
	flag := false

	for _, v := range searchButtons[currentSearchScreen] {
		if button == v {
			flag = true
			break
		}
	}

	return flag
}

func toggleButtonCheck(button string) {
	prefix := "✅ "
	for i, v := range searchButtons[currentSearchScreen] {
		if v == button {
			if strings.Contains(button, prefix) {
				key := strings.TrimPrefix(button, prefix)
				searchButtons[currentSearchScreen][i] = key
				delete(searchCriterias, key)
			} else {
				searchButtons[currentSearchScreen][i] = prefix + button
				searchCriterias[button] = button
			}
		}
	}
}

func SendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, mainMenuDescription)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = getMainMenuMarkup()
	_, err := bot.Send(msg)
	return err
}

func getCriteria() string {
	val := ""
	for _, v := range searchCriterias {
		val = fmt.Sprintf("Введите критерий поиска <b>%s</b>", v)
		currentCriteria = v
	}
	return val
}

func getMainMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUserButton, searchUserButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchUniversityButton, searchUniversityButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(miroButton, "https://miro.com/app/board/uXjVN5NbjoM=/"),
		),
	)
}

func getSearchMenuMarkup(searchType string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range searchButtons[searchType] {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btn, btn),
		)
		rows = append(rows, row)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(applyButton, applyButton),
	})

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func getUserSearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("user")
}

func getUniversitySearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return getSearchMenuMarkup("university")
}
