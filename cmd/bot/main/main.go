package main

import (
	"context"
	"telegram-bot/internal/clients/tg"
	"telegram-bot/internal/config"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/db"
	"telegram-bot/internal/model/messages"
)

// Параметры по умолчанию (могут быть изменены через config)
var (
	connectionStringDB string // Строка подключения к базе данных.
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

	// Оборачивание в Middleware функции обработки сообщения
	tgProcessingFuncHandler := tg.HandlerFunc(tg.ProcessingMessages)

	// Инициализация телеграм клиента.
	tgClient, err := tg.New(config, tgProcessingFuncHandler)
	if err != nil {
		logger.Fatal("Ошибка инициализации ТГ-клиента:", "err", err)
	}

	// Инициализация хранилищ (подключение к базе данных).
	pool, err := dbutils.NewDBConnect(ctx, connectionStringDB)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных:", "err", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal("Ошибка пинга БД", "err", err)
	}

	// БД информации пользователей.
	userStorage := db.NewUserStorage(pool)

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
}
