package tg

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/helpers/markdown"
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

func (c *Client) SendMessageWithMarkup(message dialog.Message) error {
	text := markdown.EscapeForMarkdown(message.Text)
	msg := tgbotapi.NewMessage(message.ChatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = message.Markup
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) SendMessage(msg dialog.Message, msgType string) error {
	switch msgType {
	case "SendMessageWithMarkup":
		return c.SendMessageWithMarkup(msg)
	case "SendCards":
		for _, card := range msg.Cards {
			msg.Text = card
			c.DeferMessage(msg)
		}
		return nil
	case "SendFile":
		return c.SendFile(msg)
	case "SendMedia":
		return c.SendMedia(msg)
	case "SendMediaGroup":
		return c.SendMediaGroup(msg)
	case "SendDBDump":
		return c.SendDBDump()
	case "EditTextAndMarkup":
		return c.EditTextAndMarkup(msg)
	case "EditMarkup":
		return c.EditMarkup(msg)
	default:
		text := markdown.EscapeForMarkdown(msg.Text)
		msg := tgbotapi.NewMessage(msg.ChatID, text)
		msg.ParseMode = "MarkdownV2"
		_, err := c.client.Send(msg)
		if err != nil {
			return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
		}
		return nil
	}
}

func (c *Client) SendFile(msg dialog.Message) error {
	fileConfig := tgbotapi.NewDocument(msg.ChatID, msg.File)
	fileConfig.Caption = msg.Caption
	_, err := c.client.Send(fileConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMedia(msg dialog.Message) error {
	fileConfig := tgbotapi.NewPhoto(msg.ChatID, msg.Photo)
	fileConfig.Caption = msg.Caption
	_, err := c.client.Send(fileConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMediaGroup(msg dialog.Message) error {
	var mediaGroup []interface{}

	for i, path := range msg.FilePaths {
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(path))
		if i == 0 {
			photo.Caption = msg.Caption
		}
		mediaGroup = append(mediaGroup, photo)
	}

	mg := tgbotapi.NewMediaGroup(msg.ChatID, mediaGroup)

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

	msg := dialog.Message{
		ChatID:  id,
		File:    dbDump,
		Caption: "Бэкап БД",
		Type:    "SendDBDump",
	}
	c.DeferMessage(msg)

	// TODO обрабатывать ошибку отправки файла
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
				gocron.NewAtTime(17, 0, 50),
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

func (c *Client) DeferMessage(msg dialog.Message) {
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

	for range timer.C {
		var cases []reflect.SelectCase

		for chatId, ch := range deferredMessages {
			if userCanReceiveMessage(chatId) && len(ch) > 0 {
				sc := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
				cases = append(cases, sc)
			}
		}

		if len(cases) > 0 {
			_, value, ok := reflect.Select(cases)

			if ok {
				dm := value.Interface().(dialog.Message)
				err := c.SendMessage(dm, dm.Type)
				if err != nil {
					errMsg := err.Error()
					if errMsg == "Forbidden: bot was blocked by the user" ||
						errMsg == "Forbidden: user is deactivated" {
						logger.Info(fmt.Sprintf(
							"Пользователь [%s | @%s] заблокировал бота или был деактивирован",
							dm.FirstName,
							dm.UserName,
						))
					} else {
						logger.Error("Ошибка при обработке сообщения:", "ERROR", errMsg)
					}
				}
				lastMessageTimes[dm.ChatID] = time.Now().UnixNano()
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
	if update.Message != nil {
		logger.Info(
			fmt.Sprintf(
				"[@%s | %v] %s",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Text),
		)
		msgModel.HandleMessage(dialog.Message{
			Text:           update.Message.Text,
			ChatID:         update.Message.Chat.ID,
			IsCommand:      update.Message.IsCommand(),
			BotName:        c.client.Self.UserName,
			FirstName:      update.Message.From.FirstName,
			NewChatMembers: update.Message.NewChatMembers,
			LeftChatMember: update.Message.LeftChatMember,
		})
	} else if update.CallbackQuery != nil {
		logger.Info(fmt.Sprintf("[@%s][%v] Callback: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data))
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := c.client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "ERROR", err)
		}
		msgModel.HandleButton(dialog.Message{
			CallbackQuery: update.CallbackQuery,
			MsgID:         update.CallbackQuery.Message.MessageID,
			ChatID:        update.CallbackQuery.Message.Chat.ID,
		})
	}
}

// isDuplicateEdit Проверка на действия, которые не приведут к изменениям
func isDuplicateEdit(
	msg dialog.Message,
	onlyMarkup bool,
) bool {
	newText := msg.NewText
	oldText := msg.Text
	newMarkup := msg.Markup
	oldMarkup := msg.CallbackQuery.Message.ReplyMarkup

	if onlyMarkup {
		if &newMarkup == oldMarkup {
			return true
		}
	} else {
		if (&newMarkup == oldMarkup) && (newText == oldText) {
			return true
		}
	}

	return false
}

// EditTextAndMarkup Замена текста и инлайн-кнопок.
// Их нажатие ожидает коллбек-ответ.
func (c *Client) EditTextAndMarkup(
	msg dialog.Message,
) error {
	if !isDuplicateEdit(msg, false) {
		chatID := msg.ChatID
		msgID := msg.MsgID
		text := markdown.EscapeForMarkdown(msg.NewText)

		msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, msg.Markup)
		msg.ParseMode = "MarkdownV2"
		_, err := c.client.Send(msg)
		if err != nil {
			logger.Error("Ошибка при редактировании текста и кнопок сообщения", "ERROR", err)
			return errors.Wrap(err, "client.Send with text and inline-buttons edit")
		}
	}
	return nil
}

// EditMarkup Замена инлайн-кнопок
func (c *Client) EditMarkup(msg dialog.Message) error {
	if !isDuplicateEdit(msg, true) {
		chatID := msg.ChatID
		msgID := msg.MsgID
		markup := msg.Markup

		_msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, markup)
		_, err := c.client.Send(_msg)
		if err != nil {
			logger.Error("Ошибка при редактировании текста и кнопок сообщения", "ERROR", err)
			return errors.Wrap(err, "client.Send with text and inline-buttons edit")
		}
	}
	return nil
}
