package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/card"
	"telegram-bot/internal/model/db"
	"telegram-bot/internal/model/search"
	"time"

	route "telegram-bot/internal/server"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

// Model Модель бота (клиент, хранилище, поиск)
type Model struct {
	client  *gotgbot.Bot       // Клиент бота.
	storage db.UserDataStorage // Хранилище пользовательской информации.
	search  search.Engine      // Поиск.
}

// This bot shows how to use this library to server a webapp.
// Webapps are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
// For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
// Simply install ngrok, make an account on the website, and run:
// `ngrok http 8080`
// Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
// from the samples/webappBot directory:
// `URL="<your_url_here>" TOKEN="<your_token_here>" go run .`
// Then, simply send /start to your bot, and enjoy your webapp demo.
//
// This example also demonstrates how to use the updater's handler in a user-provided server.
func main() {
	// Get token from the environment variable
	// token := os.Getenv("TOKEN")
	token := "6587208797:AAHfjqzK8moFOdhTPUtklCFtV6dRphpiKBc"
	if token == "" {
		logger.Fatal("TOKEN environment variable is empty")
	}

	// This MUST be an HTTPS URL for telegram to accept it.
	// webappURL := os.Getenv("URL")
	// webappURL := "https://t.me/LigaCifry_bot/Liga_Cifry"
	webappURL := "https://3de7-46-138-62-113.ngrok-free.app"
	if webappURL == "" {
		logger.Fatal("URL environment variable is empty")
	}

	// Get the webhook secret from the environment variable.
	// webhookSecret := os.Getenv("WEBHOOK_SECRET")
	webhookSecret := "2fPu8S36HBJhnifQV8ehc7Psruj_6oxbixZch5NZuygDbBpTe"
	if webhookSecret == "" {
		logger.Fatal("WEBHOOK_SECRET environment variable is empty")
	}

	pool, err := dbutils.NewDBConnect(
		context.TODO(),
		5,
		"postgresql://postgres:7406@localhost:5432/liga_cifry",
	)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных:", "err", err)
	}

	// БД информации пользователей.
	userStorage := db.NewUserStorage(pool)

	// Механизм поиска.
	searchEngine := search.Init()

	// Create our bot.
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		logger.Fatal("failed to create new bot: ", "err", err)
	}

	model := &Model{
		client:  b,
		search:  searchEngine,
		storage: userStorage,
	}

	// Create updater and dispatcher to handle updates in a simple manner.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Error("an error occurred while handling update:", "err", err)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot and send the URL
	dispatcher.AddHandler(
		handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
			// We can wrap commands with anonymous functions to pass in extra variables, like the webapp URL, or other
			// configuration.
			return start(b, ctx, webappURL)
		}),
	)

	dispatcher.AddHandler(
		handlers.NewCommand("allpersonalcards", func(b *gotgbot.Bot, ctx *ext.Context) error {
			return allPersonalCards(model, ctx)
		}),
	)
	// Answer callback query sent in the /start command.
	// FIX repeatable
	dispatcher.AddHandler(
		handlers.NewCallback(
			callbackquery.Equal("back_to_menu"),
			func(b *gotgbot.Bot, ctx *ext.Context) error {
				return openMenu(b, ctx, webappURL)
			},
		),
	)

	// We add the bot webhook to our updater, such that we can populate the updater's http.Handler.
	err = updater.AddWebhook(b, b.Token, &ext.AddWebhookOpts{SecretToken: webhookSecret})
	if err != nil {
		logger.Fatal("Failed to add bot webhooks to updater: " + err.Error())
	}

	// We select a subpath to specify where the updater handler is found on the http.Server.
	updaterSubpath := "/bots/"
	err = updater.SetAllBotWebhooks(webappURL+updaterSubpath, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookSecret,
	})
	if err != nil {
		logger.Fatal("Failed to set bot webhooks: " + err.Error())
	}

	// Setup new HTTP server mux to handle different paths.
	mux := http.NewServeMux()
	// This serves the hcome page.
	mux.HandleFunc("/", route.Index(webappURL))
	// This serves our "validation" API, which checks if the input data is valid.
	mux.HandleFunc("/validate", route.Validate(token))
	// This serves the updater's webhook handler.
	mux.HandleFunc(updaterSubpath, updater.GetHandlerFunc(updaterSubpath))

	server := http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8080",
	}

	log.Printf("%s has been started...\n", b.User.Username)
	// Start the webserver displaying the page.
	// Note: ListenAndServe is a blocking operation, so we don't need to call updater.Idle() here.
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("failed to listen and serve: " + err.Error())
	}
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context, webappURL string) error {
	_, err := b.SendMessage(
		ctx.Message.From.Id,
		fmt.Sprintf(
			"👋 Привет, <b>%v</b>.\nМогу помочь найти карточку компетенций или организации",
			ctx.Message.From.FirstName,
		),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{
						Text:   "🔍 Поиск индивидуальных карточек",
						WebApp: &gotgbot.WebAppInfo{Url: webappURL},
					},
				}, {
					{
						Text:   "🔍 Поиск карточек организаций",
						WebApp: &gotgbot.WebAppInfo{Url: webappURL},
					},
				}},
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func allPersonalCards(m *Model, ctx *ext.Context) error {
	_ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rawCards, err := m.storage.ShowAllPersonalCards(_ctx)
	if err != nil {
		logger.Error("collect cards", "err", err)
	}
	cards := card.FormatPersonCards(rawCards)
	for _, card := range cards {
		_, err := m.client.SendMessage(ctx.Message.Chat.Id, card, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "↩️ Меню",
							CallbackData: "back_to_menu",
						},
					},
				},
			},
		})
		if err != nil {
			logger.Error("sending all personal cards", "err", err)
		}
	}
	return nil
}

// openMenu sends main menu.
func openMenu(b *gotgbot.Bot, ctx *ext.Context, webappURL string) error {
	cb := ctx.Update.CallbackQuery
	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{})
	if err != nil {
		return fmt.Errorf("failed to answer start callback query: %w", err)
	}

	_, err = b.SendMessage(
		cb.Message.GetChat().Id,
		fmt.Sprintf(
			"👋 Привет, <b>%v</b>.\nМогу помочь найти карточку компетенций или организации",
			cb.From.FirstName,
		),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{
						Text:   "🔍 Поиск индивидуальных карточек",
						WebApp: &gotgbot.WebAppInfo{Url: webappURL},
					},
				}, {
					{
						Text:   "🔍 Поиск карточек организаций",
						WebApp: &gotgbot.WebAppInfo{Url: webappURL},
					},
				}},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send main menu: %w", err)
	}
	return nil
}
