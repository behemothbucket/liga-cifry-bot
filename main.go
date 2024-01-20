package main

import (
	"bufio"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
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

	// Button texts
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

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI("6587208797:AAEOA1DGftSvb8S8EXqKrrWRvf_BVtWQP8o")
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Сервер запущен. Нажмите Enter для остановки...")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
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
		log.Print(searchCriterias)
	} else {
		err = SendMenu(message.Chat.ID)
	}

	if err != nil {
		log.Printf("Ошибка: %s", err.Error())
	}
}

// When we get a command, we react accordingly
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
		//log.Print(searchCriterias)
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
				//log.Print(searchCriterias)
			} else {
				searchButtons[currentSearchScreen][i] = prefix + button
				searchCriterias[button] = button
				//log.Print(searchCriterias)
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

//func collectSearchData() {
//	for true {
//		if len(searchCriterias) == 0 {
//			break
//		}
//
//		SendMenu()
//
//	}
//}

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

func getUserSearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][0], searchButtons["user"][0]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][1], searchButtons["user"][1]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][2], searchButtons["user"][2]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][3], searchButtons["user"][3]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][4], searchButtons["user"][4]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][5], searchButtons["user"][5]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["user"][6], searchButtons["user"][6]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(applyButton, applyButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
	)
}

func getUniversitySearchMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][0], searchButtons["university"][0]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][1], searchButtons["university"][1]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][2], searchButtons["university"][2]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][3], searchButtons["university"][3]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][4], searchButtons["university"][4]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][5], searchButtons["university"][5]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][6], searchButtons["university"][6]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(searchButtons["university"][7], searchButtons["university"][7]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(applyButton, applyButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
	)
}
