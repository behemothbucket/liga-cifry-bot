package main

import (
	"context"
	"telegram-bot/internal/clients/tg"
	"telegram-bot/internal/config"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/db"
	"telegram-bot/internal/model/messages"
	"telegram-bot/internal/tracing"
)

// Параметры по умолчанию (могут быть изменены через config)
var (
	connectionStringDB string // Строка подключения к базе данных.
	maxAttempts        int    // Максимальное количество попыток подключения.
)

func main() {
	logger.Info("Старт приложения")

	ctx := context.Background()

	config, err := config.New()
	if err != nil {
		logger.Fatal("Ошибка получения файла конфигурации:", "err", err)
	}

	// Изменение параметров по умолчанию из заданной конфигурации.
	setConfigSettings(config.GetConfig())

	// Оборачивание в Middleware функции обработки сообщения для метрик и трейсинга.
	tgProcessingFuncHandler := tg.HandlerFunc(tg.ProcessingMessages)
	tgProcessingFuncHandler = tracing.TracingMiddleware(tgProcessingFuncHandler)

	// Инициализация телеграм клиента.
	tgClient, err := tg.New(config, tgProcessingFuncHandler)
	if err != nil {
		logger.Fatal("Ошибка инициализации ТГ-клиента:", "err", err)
	}

	// Инициализация хранилищ (подключение к базе данных).
	dbpool, err := dbutils.NewDBConnect(ctx, maxAttempts, connectionStringDB)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных:", "err", err)
	}
	defer dbpool.Close()

	// БД информации пользователей.
	userStorage := db.NewUserStorage(dbpool)

	msgModel := messages.New(ctx, tgClient, userStorage)

	// Запуск ТГ-клиента.
	tgClient.ListenUpdates(msgModel)

	logger.Info("Завершение приложения")
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
