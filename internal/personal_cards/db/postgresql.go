package db

import (
	"context"
	"fmt"

	"telegram-bot/internal/personal_cards"
	"telegram-bot/pkg/client/postgresql"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	client postgresql.Client
	// TODO Сделать logger
}

func NewRepository(client postgresql.Client) personal_cards.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) ShowAllPersonalCards(ctx context.Context) (pc []personal_cards.PersonalCard, err error) {
	q := `SELECT * FROM personal_cards;`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	_cards, err := pgx.CollectRows(rows, pgx.RowToStructByName[personal_cards.PersonalCard])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return
	}

	return _cards, nil
}
