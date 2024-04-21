package tgfile

//
// import (
// 	"telegram-bot/internal/logger"
//
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )
//
// func CreateDocument(path string) (*tgbotapi.FileReader, error) {
// 	logger.Info("Создаю документ для отправки в Телеграм...")
//
// 	file, reader, err := tgbotapi.FilePath(path).UploadData()
// 	if err != nil {
// 		logger.Error("Ошибка при конфигурации файла для FileReader", "ERROR", err)
// 	}
//
// 	document := &tgbotapi.FileReader{
// 		Name:   file,
// 		Reader: reader,
// 	}
//
// 	logger.Info("Документ успешно создан")
//
// 	return document, nil
// }
