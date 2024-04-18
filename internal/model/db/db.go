package db

import (
	"context"
	"errors"
	"fmt"
	"telegram-bot/internal/helpers/dbutils"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/card"
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
	FindPersonCards(
		ctx context.Context,
		table string,
		data []string,
		crterions []string,
	) ([]card.PersonCard, error)
	FindOrganizationCards(
		ctx context.Context,
		table string,
		data []string,
		crterions []string,
	) ([]card.OrganizationCard, error)
	ShowAllPersonalCards(ctx context.Context) (pc []card.PersonCard, err error)
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
				"[SQL: JoinGroup] Пользователь @%s обновил свои данные",
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
		logger.Info(fmt.Sprintf("Пользователь @%s присоединился к группе", u.UserName))
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
		fmt.Sprintf("Пользователь @%s вышел из группы или был исключен", u.UserName),
	)

	return nil
}

func (s *UserStorage) FindPersonCards(
	ctx context.Context,
	table string,
	data []string,
	criterions []string,
) ([]card.PersonCard, error) {
	if len(criterions) == 0 {
		return nil, errors.New("at least one criterion is required")
	}

	var query string
	var args []interface{}

	if len(criterions) == 1 {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s ILIKE $1", table, criterions[0])
		args = append(args, "%"+data[0]+"%")
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE ", table)
		for i, criterion := range criterions {
			if i > 0 {
				query += " AND "
			}
			query += fmt.Sprintf("%s ILIKE $%d", criterion, i+1)
			args = append(args, "%"+data[i]+"%")
		}
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logger.Info(fmt.Sprintf("\nQuery:\n%s\nData:\n%v", dbutils.FormatQuery(query), data))

	var cards []card.PersonCard

	cards, err = pgx.CollectRows(rows, pgx.RowToStructByName[card.PersonCard])
	if err != nil {
		logger.Error("Ошибка при конвертации строк в структуру", "ERROR", err)
		return nil, err
	}

	return cards, nil
}

func (s *UserStorage) FindOrganizationCards(
	ctx context.Context,
	table string,
	data []string,
	criterions []string,
) ([]card.OrganizationCard, error) {
	if len(criterions) == 0 {
		return nil, errors.New("at least one criterion is required")
	}

	var query string
	var args []interface{}

	if len(criterions) == 1 {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s ILIKE $1", table, criterions[0])
		args = append(args, "%"+data[0]+"%")
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE ", table)
		for i, criterion := range criterions {
			if i > 0 {
				query += " AND "
			}
			query += fmt.Sprintf("%s ILIKE $%d", criterion, i+1)
			args = append(args, "%"+data[i]+"%")
		}
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logger.Info(fmt.Sprintf("\nQuery:\n%s\nData:\n%v", dbutils.FormatQuery(query), data))

	var cards []card.OrganizationCard

	cards, err = pgx.CollectRows(rows, pgx.RowToStructByName[card.OrganizationCard])
	if err != nil {
		logger.Error("Ошибка при конвертации строк в структуру", "ERROR", err)
		return nil, err
	}

	return cards, nil
}

func (s *UserStorage) ShowAllPersonalCards(
	ctx context.Context,
) ([]card.PersonCard, error) {
	query := "SELECT * FROM personal_cards"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logger.Info(fmt.Sprintf("\nQuery:\n%s", dbutils.FormatQuery(query)))

	var cards []card.PersonCard

	cards, err = pgx.CollectRows(rows, pgx.RowToStructByName[card.PersonCard])
	if err != nil {
		logger.Error("Ошибка при конвертации строк в структуру", "ERROR", err)
		return nil, err
	}

	return cards, nil
}
