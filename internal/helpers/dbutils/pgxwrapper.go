package dbutils

import (
	"context"
	"log"
	"strings"
	"telegram-bot/internal/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBConnect Инициализация подключения к базе данных по заданным параметрам.
func NewDBConnect(
	ctx context.Context,
	maxAttempts int,
	connString string,
) (pool *pgxpool.Pool, err error) {
	err = doWithTries(
		func() error {
			_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			pool, err = pgxpool.New(_ctx, connString)
			if err != nil {
				log.Panicf("Unable to create connection pool: %v\n", err)
				return err
			}

			return nil
		}, maxAttempts, 5*time.Second)

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal("Ошибка пинга БД", "ERROR", err)
	}

	return pool, nil
}

func doWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}

	return nil
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
