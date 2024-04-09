package dbutils

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBConnect Инициализация подключения к базе данных по заданным параметрам.
func NewDBConnect(
	ctx context.Context,
	maxAttempts int,
	connString string,
) (dbpool *pgxpool.Pool, err error) {
	err = doWithTries(
		func() error {
			_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			dbpool, err = pgxpool.New(_ctx, connString)
			if err != nil {
				log.Panicf("Unable to create connection pool: %v\n", err)
				return err
			}

			return nil
		}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Panic("Error with retires")
	}

	return dbpool, nil
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

	return
}
