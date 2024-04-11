package tg

import (
	"context"
	"fmt"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/dialog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type HandlerFunc func(tgUpdate tgbotapi.Update, c *Client, msgModel *dialog.Model)

func (f HandlerFunc) RunFunc(tgUpdate tgbotapi.Update, c *Client, msgModel *dialog.Model) {
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
	msg.ParseMode = "MarkdownV2"
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) SendMessageWithMarkup(
	text string,
	chatID int64,
	markup tgbotapi.InlineKeyboardMarkup,
) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = markup
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *dialog.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("Start listening for tg dialog")

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			// logger.Debug(fmt.Sprintf("%v", update))
			c.handlerProcessingFunc.RunFunc(update, c, msgModel)
		}
	}

	// for update := range updates {
	// 	// Функция обработки сообщений (обернутая в middleware).
	// 	c.handlerProcessingFunc.RunFunc(update, c, msgModel)
	// 	// вместо ProcessingMessages(update, c, msgModel)
	// }
}

// ProcessingMessages функция обработки сообщений.
func ProcessingMessages(update tgbotapi.Update, c *Client, msgModel *dialog.Model) {
	if update.Message != nil {
		// Пользователь написал текстовое сообщение.
		logger.Info(
			fmt.Sprintf(
				"[@%s][%v] %s",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Text,
			),
		)

		err := msgModel.HandleMessage(dialog.Message{
			Text:            update.Message.Text,
			ChatID:          update.Message.Chat.ID,
			MsgID:           update.Message.MessageID,
			FirstName:       update.Message.From.FirstName,
			NewChatMembers:  update.Message.NewChatMembers,
			LeftChatMembers: update.Message.LeftChatMember,
		})
		if err != nil {
			logger.Error("error processing message:", "err", err)
		}
	} else if update.CallbackQuery != nil {
		// Пользователь нажал кнопку.
		logger.Info(fmt.Sprintf("[@%s][%v] Callback: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data))

		// Text if not specified, nothing will be shown to the user
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

		if _, err := c.client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "err", err)
		}

		err := msgModel.HandleButton(dialog.Message{
			ChatID:        update.CallbackQuery.Message.Chat.ID,
			MsgID:         update.CallbackQuery.Message.MessageID,
			CallbackQuery: update.CallbackQuery,
		})
		if err != nil {
			logger.Error("error handle button from callback:", "err", err)
		}
	}
}

// ShowInlineButtons Отображение кнопок меню под сообщением с ответом.
// Их нажатие ожидает коллбек-ответ.
func (c *Client) ShowInlineButtons(
	chatID int64,
	msgID int,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) error {
	msg := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, markup)
	msg.ParseMode = "MarkdownV2"
	_, err := c.client.Send(msg)
	if err != nil {
		logger.Error("Ошибка отправки сообщения", "err", err)
		return errors.Wrap(err, "client.Send with inline-buttons")
	}
	return nil
}
