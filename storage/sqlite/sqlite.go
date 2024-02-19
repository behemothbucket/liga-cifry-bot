package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage struct {
	db *sql.DB
}

// New is used to create new Storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

// AddUser is used to add User from database.
func (s *Storage) AddUser(ctx context.Context, u *tgbotapi.User) error {
	q := `INSERT INTO users (id, user_name, first_name, last_name, is_bot, date_joined) VALUES (?, ?, ?, ?, ?, ?);`

	currentTime := time.Now().Format("01-02-2006 15:04:05")

	if _, err := s.db.ExecContext(ctx, q, u.ID, u.UserName, u.FirstName, u.LastName, u.IsBot, currentTime); err != nil {
		return fmt.Errorf("can't save user: %w", err)
	}

	return nil
}

// DeleteUser is used to delete User from database.
func (s *Storage) DeleteUser(ctx context.Context, u *tgbotapi.User) error {
	q := `DELETE FROM users WHERE id = ? AND userName = ?;`

	if _, err := s.db.ExecContext(ctx, q, u.ID, u.UserName); err != nil {
		return fmt.Errorf("can't delete user: %w", err)
	}

	return nil
}

// IsExist is used to check if a User exists in the database.
func (s *Storage) IsExist(ctx context.Context, u tgbotapi.User) (bool, error) {
	q := `SELECT COUNT(*) FROM users WHERE id = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, u.ID, u.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if user exist: %w", err)
	}

	return count > 0, nil
}

// Init is used to create table users.
func (s *Storage) Init(ctx context.Context) error {
	q := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		user_name TEXT,
		first_name TEXT,
		last_name TEXT,
		is_bot INTEGER,
		date_joined TEXT
	);`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
