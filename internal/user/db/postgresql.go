package db

import (
	"context"
	"log"
	"time"

	"telegram-bot/internal/user"
	"telegram-bot/pkg/client/postgresql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type repository struct {
	client postgresql.Client
	// TODO Сделать logger
}

func NewRepository(client postgresql.Client) user.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) getCurrentTime() time.Time {
	return time.Now()
}

func (r *repository) doesUserExist(ctx context.Context, userID int64) (bool, error) {
	checkUserQuery := `
    SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)
  `
	var exist bool
	err := r.client.QueryRow(ctx, checkUserQuery, userID).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *repository) JoinGroup(ctx context.Context, u *tgbotapi.User) error {
	exist, err := r.doesUserExist(ctx, u.ID)
	if err != nil {
		return err
	}

	currentTime := r.getCurrentTime()

	if exist {
		if _, err := r.client.Exec(ctx, `
      UPDATE users
      SET
        is_joined = $1, date_joined = $2, date_left = $3
      WHERE
        user_id = $4
    `, true, currentTime, nil, u.ID); err != nil {
			return err
		}
		log.Printf("[SQL: JoinGroup] Пользователь с никнеймом @%s обновил свои данные", u.UserName)
	} else {
		if _, err := r.client.Exec(ctx, `
      INSERT INTO users
        (user_id, user_name, first_name, last_name, is_bot, is_joined, date_joined, date_left)
      VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8)
    `, u.ID, u.UserName, u.FirstName, u.LastName, u.IsBot, true, currentTime, nil); err != nil {
			return err
		}
		log.Printf("Пользователь с никнеймом @%s присоединился к группе", u.UserName)
	}

	return nil
}

func (r *repository) LeaveGroup(ctx context.Context, u *tgbotapi.User) error {
	exist, err := r.doesUserExist(ctx, u.ID)
	if err != nil {
		return err
	}

	if !exist {
		log.Printf("[SQL: LeaveGroup] Пользователь с ID %d не найден", u.ID)
		return nil
	}

	q := `
    UPDATE users
    SET
      is_joined = $1, date_left = $2
    WHERE
      user_id = $3

    `

	currentTime := r.getCurrentTime()
	if _, err := r.client.Exec(ctx, q, false, currentTime, u.ID); err != nil {
		return err
	}

	log.Printf("Пользователь с никнеймом @%s вышел из группы или был исключен", u.UserName)

	return nil
}
