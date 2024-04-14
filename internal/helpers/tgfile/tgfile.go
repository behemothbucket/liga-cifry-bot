package tgfile

import (
	"os"
	"telegram-bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateDocument(filePath string) (*tgbotapi.FileBytes, error) {
	logger.Info("Создаю документ для отправки в Телеграм...")
	file, err := os.Open(filePath)
	if err != nil {
		return &tgbotapi.FileBytes{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return &tgbotapi.FileBytes{}, err
	}
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		return &tgbotapi.FileBytes{}, err
	}

	document := tgbotapi.FileBytes{Name: fileInfo.
		Name(), Bytes: fileBytes}

	logger.Info("Документ успешно создан")

	return &document, nil
}
