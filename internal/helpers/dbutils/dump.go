package dbutils

import (
	"fmt"
	"os/exec"
	"telegram-bot/internal/logger"
	"time"
)

func CreateDBDump() (string, error) {
	currentTime := time.Now().Format("02_01_2006_15_04_05")
	filePath := fmt.Sprintf("./dumps/backup_%s.sql", currentTime)

	logger.Info("Создаю файл дампа БД...")

	cmdStr := fmt.Sprintf("/usr/bin/pg_dump liga_cifry > %s", filePath)
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		logger.Error("Ошибка при выполнении дампа", "ERROR", err)
		return "", err
	}

	logger.Info("Дамп успешно создан")

	return filePath, nil
}
