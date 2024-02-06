package storage

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	AddUser(ctx context.Context, u *tgbotapi.User) error
	DeleteUser(ctx context.Context, u *tgbotapi.User) error
	IsExist(ctx context.Context, u *tgbotapi.User) (bool, error)
}
