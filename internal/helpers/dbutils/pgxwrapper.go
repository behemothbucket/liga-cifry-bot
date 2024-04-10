package dbutils

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool   *pgxpool.Pool
	pgOnce sync.Once
)

// NewDBConnect Инициализация подключения к базе данных по заданным параметрам.
func NewDBConnect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var err error

	pgOnce.Do(func() {
		db, _err := pgxpool.New(ctx, connString)
		if _err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", _err)
			return
		}

		pool = db
	})

	if err != nil {
		return nil, err
	}

	return pool, nil
}
