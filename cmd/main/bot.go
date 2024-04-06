package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"telegram-bot/internal/config"
	"telegram-bot/internal/personal_cards"
	PC "telegram-bot/internal/personal_cards/db"
	"telegram-bot/pkg/client/postgresql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SQLTEST() []personal_cards.PersonalCard {
	cfg := config.GetConfig()
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, *cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	repository := PC.NewRepository(postgreSQLClient)

	cards, err := repository.ShowAllPersonalCards(context.TODO())
	if err != nil {
		log.Fatalf("%v", err)
	}

	return cards
}

type Bot struct {
	bot                      *tgbotapi.BotAPI
	TelegramApiToken         string
	searchMode               bool
	currentSearchScreen      string
	userSearchCriteria       map[string]string
	universitySearchCriteria map[string]string
}

func newBot() *Bot {
	token := getBotToken()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &Bot{
		bot:                      bot,
		TelegramApiToken:         token,
		userSearchCriteria:       map[string]string{},
		universitySearchCriteria: map[string]string{},
	}
}

func getBotToken() string {
	token, exists := os.LookupEnv("TELEGRAM_BOT_TOKEN")

	if !exists {
		log.Panic("Token not found in .env")
	}

	return token
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.handleUpdate(ctx, update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.handleMessage(ctx, update.Message)
	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery, ctx)
	}
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if message.From == nil {
		return
	}

	logMessage(message)

	switch {
	case message.IsCommand() && getChatType(message) == "private":
		b.handleCommand(message)
	case b.searchMode:
		b.sendAcceptMessage(message)
		b.sendLoadMoreMessage(message)
	case !isValidMessageText(message) && getChatType(message) == "private":
		b.sendMainMenu(message)
	case getChatType(message) == "private":
		b.sendMainMenu(message)
		// default:
		// handleIfSubscriptionEvent(ctx, message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	command := message.Text
	botName := fmt.Sprintf("@%s", b.bot.Self.UserName)

	switch command {
	case "/menu", "/menu" + botName, "/start", "/start" + botName:
		b.sendMainMenu(message)
	}
}

func getChatType(message *tgbotapi.Message) string {
	return message.Chat.Type
}

func (b *Bot) handleButton(query *tgbotapi.CallbackQuery, ctx context.Context) {
	var text string

	markup := mainMenuMarkup
	message := query.Message

	switch query.Data {
	case searchUserButton:
		b.currentSearchScreen = "user"
		text = searchMenuDescription
		markup = b.getCurrentSearchMarkup()
	case searchUniversityButton:
		b.currentSearchScreen = "university"
		text = searchMenuDescription
		markup = b.getCurrentSearchMarkup()
	case backButton:
		text = mainMenuDescription
		markup = mainMenuMarkup
	case menuButton:
		text = mainMenuDescription
		// resetCriteriaButtons() // TODO сбрасывать кнопки и чистить критерии после найденной карточки, а не здесь
		b.sendMainMenu(message)
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		return
	case applyButton:
		var criteria map[string]string

		if b.currentSearchScreen == "user" {
			criteria = b.userSearchCriteria
		} else {
			criteria = b.universitySearchCriteria
		}

		if len(criteria) == 0 {
			text = "️❗️Пожалуйста, выберите хотя-бы один критерий поиска."
			markup = b.getCurrentSearchMarkup()
		} else {
			text = b.getCriterion()
			markup = getCancelMenuMarkup()
			b.searchMode = true
		}

	case cancelSearchButton:
		b.resetCriteriaButtons()
		b.searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		b.sendMainMenu(message)
		return

	case b.criterionButtonIsClicked(query.Data):
		b.toggleCriterionButton(query.Data)
		text = searchMenuDescription
		markup = b.getCurrentSearchMarkup()
	// case printFirstPersonalCard:
	// 	card := b.SpreadsheetConfig.getCardByNumber(b.SpreadsheetConfig.personalSheetTitle, 1)
	// 	b.SendMessage(Message{
	// 		chatID:      message.Chat.ID,
	// 		text:        card,
	// 		groupName:   message.Chat.Type,
	// 		replyMarkup: &backToMainMenuMarkup,
	// 		parseMode:   tgbotapi.ModeHTML,
	// 	})
	// 	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	// 	b.bot.Send(callbackCfg)
	case printAllPersonalCards:
		cards := SQLTEST()
		log.Printf("%s хочет получить все карточки пользователей", message.Chat.UserName)
		for _, card := range cards {
			formattedText := fmt.Sprintf(
				`<b>🧑‍💼ФИО</b>
%s

<b>📍Город</b>
%s

<b>🏛Организация</b>
%s

<b>🤝Должность</b>
%s

<b>📝Экспертные компетенции</b>
%s

<b>🤝Направления возможного сотрудничества</b>
%s

<b>📱Контакты для связи</b>
%s`,
				card.Fio, card.City, card.Organization, card.Job_title, card.Expert_competencies, card.Possible_cooperation, card.Contacts,
			)
			b.SendMessage(Message{
				chatID:      message.Chat.ID,
				text:        formattedText,
				groupName:   message.Chat.Type,
				replyMarkup: &backToMainMenuMarkup,
				parseMode:   tgbotapi.ModeHTML,
			})
		}
		b.sendLoadMoreMessage(message)
		// b.SendMessage(Message{
		// 	chatID:      message.Chat.ID,
		// 	text:        "✅ Показаны все персональные карточки",
		// 	groupName:   message.Chat.Type,
		// 	replyMarkup: nil,
		// 	parseMode:   tgbotapi.ModeHTML,
		// })
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		// case printFirstOrganizationCard:
		// 	card := b.SpreadsheetConfig.getCardByNumber(b.SpreadsheetConfig.organizationSheetTitle, 1)
		// 	b.SendMessage(Message{
		// 		chatID:      message.Chat.ID,
		// 		text:        card,
		// 		groupName:   message.Chat.Type,
		// 		replyMarkup: &backToMainMenuMarkup,
		// 		parseMode:   tgbotapi.ModeHTML,
		// 	})
		// 	callbackCfg := tgbotapi.NewCallback(query.ID, "")
		// 	b.bot.Send(callbackCfg)

	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.bot.Send(msg)
}

func (b *Bot) getCurrentSearchMarkup() tgbotapi.InlineKeyboardMarkup {
	if b.currentSearchScreen == "user" {
		return getSearchMenuMarkup("user")
	} else {
		return getSearchMenuMarkup("university")
	}
}
