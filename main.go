package main

import (
	"bufio"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var bot *tgbotapi.BotAPI

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
		handleMessage(update)
		break
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}
