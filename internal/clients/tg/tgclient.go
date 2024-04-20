package tg

import (
	"context"
	"fmt"
	"os"
	"sync"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/helpers/tgfile"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/dialog"
	"time"

	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

const (
	sendCooldownPerUser = int64(time.Second / 3)
	sendInterval        = time.Second
	deferMessages       = false
)

var (
	deferredMessages = make(map[int64]chan dialog.Message)
	lastMessageTimes = make(map[int64]int64)
)

type Client struct {
	sync.RWMutex
	client *tgbotapi.BotAPI
}

type TokenGetter interface {
	Token() string
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) SendMessageWithMarkup(
	chatID int64,
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = markup
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) SendKeyboard(
	chatID int64,
	text string,
	markup *tgbotapi.ReplyKeyboardMarkup,
) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = markup
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) SendCards(chatID int64, cards []string) error {
	for _, card := range cards {
		msg := dialog.Message{
			ChatID: chatID,
			Text:   card,
			Markup: &dialog.MarkupCardMenu,
		}
		c.DeferMessageWithMarkup(msg)
	}
	return nil
}

func (c *Client) SendFile(chatID int64, file *tgbotapi.FileReader, caption string) error {
	fileConfig := tgbotapi.NewDocument(chatID, file)
	fileConfig.Caption = caption
	_, err := c.client.Send(fileConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMedia(chatID int64, file *tgbotapi.FileReader, caption string) error {
	fileConfig := tgbotapi.NewPhoto(chatID, file)
	fileConfig.Caption = caption
	_, err := c.client.Send(fileConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMediaGroup(chatID int64, paths []string, caption string) error {
	var mediaGroup []interface{}

	for i, path := range paths {
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(path))
		if i == 0 {
			photo.Caption = caption
		}
		mediaGroup = append(mediaGroup, photo)
	}

	mg := tgbotapi.NewMediaGroup(chatID, mediaGroup)

	_, err := c.client.SendMediaGroup(mg)
	return err
}

func (c *Client) SendDBDump() error {
	// TODO  Get IDs from ENV
	id := int64(5587823077)

	filePath, err := dbutils.CreateDBDump()
	if err != nil {
		logger.Error("Ошибка при создании дампа базы данных:", err)
	}

	dbDump, err := tgfile.CreateDocument(filePath)
	if err != nil {
		logger.Error("Ошибка при создании дампа БД", "ERROR", err)
	}

	logger.Info("Начинаю рассылку дампа...")

	err = c.SendFile(id, dbDump, "Бэкап БД")
	if err != nil {
		logger.Error("Ошибка при отправке файла в телеграм:", err)
	}

	logger.Info(fmt.Sprintf("Файл отправлен %s", filePath))

	err = os.Remove(filePath)
	if err != nil {
		logger.Error("Ошибка при удалении временного файла:", err)
	}

	logger.Info(fmt.Sprintf("Файл удален %s", filePath))

	return nil
}

func (c *Client) StartDBJob(ctx context.Context) {
	logger.Info("Старт джобы по бэкапу БД")
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Error("Ошибка в старте шедулера", "ERROR", err)
	}

	j, err := s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 3, 0),
			),
		), gocron.NewTask(c.SendDBDump),
	)
	if err != nil {
		logger.Error("Ошибка в создании джобы", "ERROR", err)
	}

	logger.Info(fmt.Sprintf("Джоба: [%s] %s", j.Name(), j.ID().String()))

	s.Start()

	<-ctx.Done()

	err = s.Shutdown()
	if err != nil {
		logger.Error("Ошибка в завершении джобы", "ERROR", err)
	}

	logger.Info("Приложение завершило работу, job")
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *dialog.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("Начинаю следить за обновлениями")

	var wg sync.WaitGroup

	for update := range updates {
		select {
		case <-ctx.Done():
			logger.Info("Приложение завершило работу")
			return
		default:
			wg.Add(1)
			go func(update tgbotapi.Update) {
				defer wg.Done()
				ProcessingMessages(update, c, msgModel)
			}(update)
		}
	}
}

func (c *Client) DeferMessageWithMarkup(msg dialog.Message) {
	chatId := msg.ChatID

	c.Lock()
	defer c.Unlock()

	if _, ok := deferredMessages[chatId]; !ok {
		deferredMessages[chatId] = make(chan dialog.Message, 100)
	}

	deferredMessages[chatId] <- msg
}

func (c *Client) SendDeferredMessages() {
	timer := time.NewTicker(sendInterval)
	defer timer.Stop()

	for range timer.C {
		for chatID, ch := range deferredMessages {
			if userCanReceiveMessage(chatID) && len(ch) > 0 {
				select {
				case dm := <-ch:
					err := c.SendMessageWithMarkup(dm.ChatID, dm.Text, dm.Markup)
					if err != nil {
						logger.Error("Ошибка в отправке отложенного сообщения", "ERROR", err)
					}
					lastMessageTimes[dm.ChatID] = time.Now().UnixNano()
				default:
					// Нет доступных сообщений для отправки
					logger.Debug("Нет сообщений для отправки...")
				}
			}
		}
	}
}

func userCanReceiveMessage(userId int64) bool {
	t, ok := lastMessageTimes[userId]

	return !ok || t+sendCooldownPerUser <= time.Now().UnixNano()
}

// ProcessingMessages Функция обработки сообщений.
func ProcessingMessages(
	update tgbotapi.Update,
	c *Client,
	msgModel *dialog.Model,
) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		logger.Info(
			fmt.Sprintf(
				"[@%s | %v] %s",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Text,
			),
		)

		err := msgModel.HandleMessage(dialog.Message{
			Text:           update.Message.Text,
			ChatID:         update.Message.Chat.ID,
			MsgID:          update.Message.MessageID,
			IsCommand:      update.Message.IsCommand(),
			BotName:        c.client.Self.UserName,
			FirstName:      update.Message.From.FirstName,
			NewChatMembers: update.Message.NewChatMembers,
			LeftChatMember: update.Message.LeftChatMember,
		})
		if err != nil {
			errMsg := err.Error()
			if errMsg == "Forbidden: bot was blocked by the user" ||
				errMsg == "Forbidden: user is deactivated" {
				logger.Info(fmt.Sprintf(
					"Пользователь [%s | @%s] заблокировал бота или был деактивирован",
					update.Message.From.FirstName,
					update.Message.From.UserName,
				))
			} else {
				logger.Error("Ошибка при обработке сообщения:", "ERROR", errMsg)
			}
		}
	} else if update.CallbackQuery != nil {
		logger.Info(fmt.Sprintf("[@%s][%v] Callback: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data))
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := c.client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "ERROR", err)
		}
		err := msgModel.HandleButton(dialog.Message{
			ChatID:        update.CallbackQuery.Message.Chat.ID,
			MsgID:         update.CallbackQuery.Message.MessageID,
			CallbackQuery: update.CallbackQuery,
			FirstName:     update.CallbackQuery.From.FirstName,
			Text:          update.CallbackQuery.Message.Text,
		})
		if err != nil {
			logger.Error("error handle button from callback:", "ERROR", err)
		}
	}
}

// isDuplicateEdit Проверка на действия, которые не приведут к изменениям
func isDuplicateEdit(
	msg dialog.Message,
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
	onlyMarkup bool,
) bool {
	oldText := msg.Text
	oldMarkup := msg.Markup

	if onlyMarkup {
		if markup == oldMarkup {
			return true
		}
	} else {
		if (markup == oldMarkup) && (text == oldText) {
			return true
		}
	}

	return false
}

// EditTextAndMarkup Замена текста и инлайн-кнопок.
// Их нажатие ожидает коллбек-ответ.
func (c *Client) EditTextAndMarkup(
	msg dialog.Message,
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) error {
	if !isDuplicateEdit(msg, text, markup, false) {
		chatID := msg.ChatID
		msgID := msg.MsgID

		msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, *markup)
		msg.ParseMode = "HTML"
		_, err := c.client.Send(msg)
		if err != nil {
			logger.Error("Ошибка при редактировании текста и кнопок сообщения", "ERROR", err)
			return errors.Wrap(err, "client.Send with text and inline-buttons edit")
		}
	}
	return nil
}

// EditMarkup Замена инлайн-кнопок
func (c *Client) EditMarkup(msg dialog.Message, markup *tgbotapi.InlineKeyboardMarkup) error {
	if !isDuplicateEdit(msg, "", markup, true) {
		chatID := msg.ChatID
		msgID := msg.MsgID

		_msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, *markup)
		_, err := c.client.Send(_msg)
		if err != nil {
			logger.Error("Ошибка при редактировании текста и кнопок сообщения", "ERROR", err)
			return errors.Wrap(err, "client.Send with text and inline-buttons edit")
		}
	}
	return nil
}
