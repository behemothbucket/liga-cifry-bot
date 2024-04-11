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
	dialog "telegram-bot/internal/model/dialog"
	"telegram-bot/internal/model/search"
)

// Параметры по умолчанию (могут быть изменены через config)
var (
	connectionStringDB string // Строка подключения к базе данных.
)

func main() {
	logger.Info("Старт приложения")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

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

	// Механизм поиска.
	searchEngine := search.Init()

	msgModel := dialog.New(ctx, tgClient, userStorage, searchEngine)

	// Запуск ТГ-клиента.
	go tgClient.ListenUpdates(ctx, msgModel)

	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		logger.Error("Ошибка в завершении работы приложения")
	}
	cancel()
	logger.Info("Завершение приложения")
}

// Замена параметров по умолчанию параметрами из конфиг.файла.
func setConfigSettings(config config.Config) {
	if config.ConnectionStringDB != "" {
		connectionStringDB = config.ConnectionStringDB
	}
}
