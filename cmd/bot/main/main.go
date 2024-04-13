package main

import (
	"bufio"
	"context"
	"os"
	"telegram-bot/internal/clients/tg"
	"telegram-bot/internal/config"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/db"
	"telegram-bot/internal/model/dialog"
	"telegram-bot/internal/model/search"
)

// Параметры по умолчанию (могут быть изменены через config)
var (
	connectionStringDB string // Строка подключения к базе данных.
	maxAttempts        int    // Максимальное количестко попыток подключения к БД
)

func main() {
	logger.Info("Старт приложения")

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	dbConfig, err := config.New()
	if err != nil {
		logger.Fatal("Ошибка получения файла конфигурации:", "err", err)
	}

	// Изменение параметров по умолчанию из заданной конфигурации.
	setConfigSettings(dbConfig.GetConfig())

	// Инициализация телеграм клиента.
	tgClient, err := tg.New(dbConfig)
	if err != nil {
		logger.Fatal("Ошибка инициализации ТГ-клиента:", "err", err)
	}

	// Инициализация хранилищ (подключение к базе данных).
	pool, err := dbutils.NewDBConnect(context.TODO(), maxAttempts, connectionStringDB)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных:", "err", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal("Ошибка пинга БД", "err", err)
	}

	// БД информации пользователей.
	userStorage := db.NewUserStorage(pool)

	// Механизм поиска.
	searchEngine := search.Init()

	// Инициализация основной модели.
	msgModel := dialog.New(ctx, tgClient, userStorage, searchEngine)

	// Pass cancellable context to goroutine
	go tgClient.ListenUpdates(ctx, msgModel)

	// Tell the user the bot is online
	logger.Info("Start listening for updates. Press enter to stop...")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

// Замена параметров по умолчанию параметрами из конфиг.файла.
func setConfigSettings(config config.Config) {
	if config.ConnectionStringDB != "" {
		connectionStringDB = config.ConnectionStringDB
	}
	if config.MaxAttempts != 0 {
		maxAttempts = config.MaxAttempts
	}
}
