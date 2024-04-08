package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"telegram-bot/internal/user"
	"telegram-bot/pkg/client/postgresql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) JoinGroup(ctx context.Context, user tgbotapi.User) error {
	q := `
		INSERT INTO users
		    (user_id, user_name, first_name, last_name, is_bot, is_joined, date_joined, date_left)
		VALUES
		       ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	log.Print("SQL Query", formatQuery(q))

	if err := r.client.QueryRow(ctx, q, author.Name, 123).Scan(&author.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) LeaveGroup(ctx context.Context) error {
	q := `SELECT * FROM public.personal_cards;`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	return nil
}
