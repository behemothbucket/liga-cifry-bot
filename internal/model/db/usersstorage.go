package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/card/person"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) UserDataStorage {
	return &UserStorage{
		db: pool,
	}
}

// UserDataStorage Интерфейс для работы с хранилищем данных.
type UserDataStorage interface {
	JoinGroup(ctx context.Context, u *tgbotapi.User) error
	LeaveGroup(ctx context.Context, u *tgbotapi.User) error
	CheckIfUserExist(ctx context.Context, userID int64) (bool, error)
	FindCards(
		ctx context.Context,
		table string,
		data []string,
		crterions []string,
	) ([]person.PersonCard, error)
	ShowAllPersonalCards(ctx context.Context) (pc []person.PersonCard, err error)
}

// CheckIfUserExist Проверка существования пользователя в базе данных.
func (s *UserStorage) CheckIfUserExist(ctx context.Context, userID int64) (bool, error) {
	// Запрос на выборку пользователя.
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`

	// Выполнение запроса на получение данных.
	var exist bool
	err := s.db.QueryRow(ctx, query, userID).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

// JoinGroup Добавление пользователя в базу данных.
func (s *UserStorage) JoinGroup(
	ctx context.Context,
	u *tgbotapi.User,
) error {
	exist, err := s.CheckIfUserExist(ctx, u.ID)
	if err != nil {
		return err
	}

	if exist {
		if _, err := s.db.Exec(ctx, `
      UPDATE users
      SET
        is_joined = $1, date_joined = $2, date_left = $3
      WHERE
        user_id = $4
    `, true, time.Now(), nil, u.ID); err != nil {
			return err
		}
		logger.Info(
			fmt.Sprintf(
				"[SQL: JoinGroup] Пользователь с никнеймом @%s обновил свои данные",
				u.UserName,
			),
		)
	} else {
		if _, err := s.db.Exec(ctx, `
      INSERT INTO users
        (user_id, user_name, first_name, last_name, is_bot, is_joined, date_joined, date_left)
      VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8)
    `, u.ID, u.UserName, u.FirstName, u.LastName, u.IsBot, true, time.Now(), nil); err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Пользователь с никнеймом @%s присоединился к группе", u.UserName))
	}

	return nil
}

func (s *UserStorage) LeaveGroup(ctx context.Context, u *tgbotapi.User) error {
	exist, err := s.CheckIfUserExist(ctx, u.ID)
	if err != nil {
		return err
	}

	if !exist {
		logger.Info(fmt.Sprintf("[SQL: LeaveGroup] Пользователь с ID %d не найден", u.ID))
		return nil
	}

	q := `
    UPDATE users
    SET
      is_joined = $1, date_left = $2
    WHERE
      user_id = $3;
    `

	if _, err := s.db.Exec(ctx, q, false, time.Now(), u.ID); err != nil {
		return err
	}

	logger.Info(
		fmt.Sprintf("Пользователь с никнеймом @%s вышел из группы или был исключен", u.UserName),
	)

	return nil
}

func (s *UserStorage) FindCards(
	ctx context.Context,
	table string,
	data []string,
	criterions []string,
) ([]person.PersonCard, error) {
	if len(criterions) == 0 {
		return nil, errors.New("at least one criterion is required")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s ILIKE $1", table, criterions[0])

	for i := range data {
		query += fmt.Sprintf(" AND %s ILIKE $%d", data[i], i+2)
	}

	args := make([]interface{}, len(criterions)+1)
	args[0] = "%" + criterions[0] + "%"
	for i := range criterions {
		args[i+1] = "%" + criterions[i] + "%"
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []person.PersonCard

	for rows.Next() {
		var card person.PersonCard
		err := rows.Scan(
			&card.ID,
			&card.Fio,
			&card.City,
			&card.Organization,
			&card.JobTitle,
			&card.ExpertCompetencies,
			&card.PossibleCooperation,
			&card.Contacts,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows found")
		}
		return nil, err
	}

	return cards, nil
}

func (s *UserStorage) ShowAllPersonalCards(
	ctx context.Context,
) (pc []person.PersonCard, err error) {
	q := `SELECT * FROM public.personal_cards;`

	rows, err := s.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := pgx.CollectRows(rows, pgx.RowToStructByName[person.PersonCard])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return
	}

	return cards, nil
}
