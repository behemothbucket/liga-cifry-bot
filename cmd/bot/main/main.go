package main

import (
	"context"
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

	ctx := context.Background()

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
	pool, err := dbutils.NewDBConnect(ctx, maxAttempts, connectionStringDB)
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

	// Запуск ТГ-клиента.
	tgClient.ListenUpdates(ctx, msgModel)
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
