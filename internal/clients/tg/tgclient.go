package tg

import (
	"fmt"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type HandlerFunc func(tgUpdate tgbotapi.Update, c *Client, msgModel *messages.Model)

func (f HandlerFunc) RunFunc(tgUpdate tgbotapi.Update, c *Client, msgModel *messages.Model) {
	f(tgUpdate, c, msgModel)
}

type Client struct {
	client                *tgbotapi.BotAPI
	handlerProcessingFunc HandlerFunc // Функция обработки входящих сообщений.
}

type TokenGetter interface {
	Token() string
}

func New(tokenGetter TokenGetter, handlerProcessingFunc HandlerFunc) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка NewBotAPI")
	}

	return &Client{
		client:                client,
		handlerProcessingFunc: handlerProcessingFunc,
	}, nil
}

func (c *Client) SendMessage(text string, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("Start listening for tg messages")

	for update := range updates {
		// Функция обработки сообщений (обернутая в middleware).
		c.handlerProcessingFunc.RunFunc(update, c, msgModel)
		// вместо ProcessingMessages(update, c, msgModel)
	}
}

// ProcessingMessages функция обработки сообщений.
func ProcessingMessages(tgUpdate tgbotapi.Update, c *Client, msgModel *messages.Model) {
	if tgUpdate.Message != nil {
		// Пользователь написал текстовое сообщение.
		logger.Info(
			fmt.Sprintf(
				"[@%s][%v] %s",
				tgUpdate.Message.From.UserName,
				tgUpdate.Message.From.ID,
				tgUpdate.Message.Text,
			),
		)
		err := msgModel.IncomingMessage(messages.Message{
			Text:            tgUpdate.Message.Text,
			ChatID:          tgUpdate.Message.Chat.ID,
			FirstName:       tgUpdate.Message.From.FirstName,
			NewChatMembers:  tgUpdate.Message.NewChatMembers,
			LeftChatMembers: tgUpdate.Message.LeftChatMember,
		})
		if err != nil {
			logger.Error("error processing message:", "err", err)
		}
	} else if tgUpdate.CallbackQuery != nil {
		// Пользователь нажал кнопку.
		logger.Info(fmt.Sprintf("[@%s][%v] Callback: %s", tgUpdate.CallbackQuery.From.UserName, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Data))
		callback := tgbotapi.NewCallback(tgUpdate.CallbackQuery.ID, tgUpdate.CallbackQuery.Data)
		if _, err := c.client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "err", err)
		}
		// if err := deleteInlineButtons(c, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Message.MessageID, tgUpdate.CallbackQuery.Message.Text); err != nil {
		// 	logger.Error("Ошибка удаления кнопок:", "err", err)
		// }
		err := msgModel.HandleButton(messages.Message{
			CallbackQuery: tgUpdate.CallbackQuery,
		})
		if err != nil {
			logger.Error("error handle button from callback:", "err", err)
		}
	}
}

// ShowInlineButtons Отображение кнопок меню под сообщением с ответом.
// Их нажатие ожидает коллбек-ответ.
func (c *Client) ShowInlineButtons(
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
	chatID int64,
) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	msg.ParseMode = "HTML"
	_, err := c.client.Send(msg)
	if err != nil {
		logger.Error("Ошибка отправки сообщения", "err", err)
		return errors.Wrap(err, "client.Send with inline-buttons")
	}
	return nil
}

func deleteInlineButtons(c *Client, chatID int64, msgID int, sourceText string) error {
	msg := tgbotapi.NewEditMessageText(chatID, msgID, sourceText)
	_, err := c.client.Send(msg)
	if err != nil {
		logger.Error("Ошибка отправки сообщения", "err", err)
		return errors.Wrap(err, "client.Send remove inline-buttons")
	}
	return nil
}
