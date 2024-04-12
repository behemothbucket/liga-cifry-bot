package db

import (
	"context"
	"fmt"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/card/person"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	FindCard(ctx context.Context, criteria string, date string) (person.PersonCard, error)
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

func (s *UserStorage) FindCard(
	ctx context.Context,
	criteria string,
	data string,
) (person.PersonCard, error) {
	query := `SELECT * FROM personal_cards WHERE fio ILIKE $1`

	var card person.PersonCard

	err := s.db.QueryRow(ctx, query, "%"+data+"%").Scan(
		&card.Fio,
		&card.City,
		&card.Organization,
		&card.Job_title,
		&card.Expert_competencies,
		&card.Possible_cooperation,
		&card.Contacts,
	)
	if err != nil {
		return person.PersonCard{Fio: ""}, err
	}

	return card, nil
}

// func (s *UserStorage) ShowAllPersonalCards(
// 	ctx context.Context,
// ) (pc []dialog.PersonalCard, err error) {
// 	q := `SELECT * FROM public.personal_cards;`
//
// 	rows, err := s.db.Query(ctx, q)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	cards, err := pgx.CollectRows(rows, pgx.RowToStructByName[dialog.PersonalCard])
// 	if err != nil {
// 		fmt.Printf("CollectRows error: %v", err)
// 		return
// 	}
//
// 	return cards, nil
// }
