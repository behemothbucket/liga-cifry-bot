package db

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/model/messages"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) messages.UserDataStorage {
	return &UserStorage{
		pool: pool,
	}
}

// CheckIfUserExist Проверка существования пользователя в базе данных.
func (s *UserStorage) CheckIfUserExist(ctx context.Context, userID int64) (bool, error) {
	// Запрос на выборку пользователя.
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`

	// Выполнение запроса на получение данных.
	var exist bool
	err := s.pool.QueryRow(ctx, query, userID).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

// JoinGroup Добавление пользователя в базу данных.
func (s *UserStorage) JoinGroup(
	ctx context.Context,
	id int64,
	userName string,
	firstName string,
	lastName string,
	isBot bool,
) error {
	exist, err := s.CheckIfUserExist(ctx, id)
	if err != nil {
		return err
	}

	if exist {
		if _, err := s.pool.Exec(ctx, `
      UPDATE users
      SET
        is_joined = $1, date_joined = $2, date_left = $3
      WHERE
        user_id = $4
    `, true, time.Now(), nil, id); err != nil {
			return err
		}
		log.Printf("[SQL: JoinGroup] Пользователь с никнеймом @%s обновил свои данные", userName)
	} else {
		if _, err := s.pool.Exec(ctx, `
      INSERT INTO users
        (user_id, user_name, first_name, last_name, is_bot, is_joined, date_joined, date_left)
      VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8)
    `, id, userName, firstName, lastName, isBot, true, time.Now(), nil); err != nil {
			return err
		}
		log.Printf("Пользователь с никнеймом @%s присоединился к группе", userName)
	}

	return nil
}

func (s *UserStorage) LeaveGroup(ctx context.Context, u *tgbotapi.User) error {
	exist, err := s.CheckIfUserExist(ctx, u.ID)
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

	if _, err := s.pool.Exec(ctx, q, false, time.Now(), u.ID); err != nil {
		return err
	}

	log.Printf("Пользователь с никнеймом @%s вышел из группы или был исключен", u.UserName)

	return nil
}

func (s *UserStorage) ShowAllPersonalCards(
	ctx context.Context,
) (pc []messages.PersonalCard, err error) {
	q := `SELECT * FROM public.personal_cards;`
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	cards, err := pgx.CollectRows(rows, pgx.RowToStructByName[messages.PersonalCard])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return
	}

	return cards, nil
}
