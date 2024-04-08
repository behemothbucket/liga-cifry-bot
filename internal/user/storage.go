package user

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repository interface {
	JoinGroup(ctx context.Context, u *tgbotapi.User) error
	LeaveGroup(ctx context.Context, u *tgbotapi.User) error
}
