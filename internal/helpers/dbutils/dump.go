package dbutils

import (
	"os/exec"
	"telegram-bot/internal/logger"
	"time"
)

func CreateDBDump() (string, error) {
	currentTime := time.Now().Format("02-01-2006_15-04")
	filePath := "./dumps/backup_" + currentTime + ".dump"

	logger.Info("Создаю файл дампа БД...")

	cmd := exec.Command("bash", "-c", "/usr/bin/pg_dump liga_cifry > "+filePath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	logger.Info("Дамп успешно создан")

	return filePath, nil
}
