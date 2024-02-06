package main

import (
	"bufio"
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

func main() {
	b := newBot()

	b.bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := b.bot.GetUpdatesChan(u)

	go b.receiveUpdates(ctx, updates)

	log.Printf("Bot [%s] starts. Press Enter to terminate...", b.bot.Self.UserName)

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}
