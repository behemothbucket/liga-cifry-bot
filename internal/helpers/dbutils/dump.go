package dbutils

import (
	"fmt"
	"os/exec"
	"telegram-bot/internal/logger"
	"time"
)

func CreateDBDump() (string, string, error) {
	currentTime := time.Now().Format("02.01.2006 15:04:05")
	filePath := "./dumps/backup.dump"

	logger.Info("Создаю файл дампа БД...")

	cmdStr := fmt.Sprintf("/usr/bin/pg_dump liga_cifry > %s", filePath)
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		logger.Debug("Ошибка при выполнении дампа", "ERROR", err)
		return "", currentTime, err
	}

	logger.Info("Дамп успешно создан")

	return filePath, currentTime, nil
}
