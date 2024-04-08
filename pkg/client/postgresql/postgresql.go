package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"telegram-bot/internal/config"
	"telegram-bot/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(
	ctx context.Context,
	maxAttempts int,
	sc config.StorageConfig,
) (dbpool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		sc.UserName,
		sc.Password,
		sc.Host,
		sc.Port,
		sc.Database,
	)
	log.Print(dsn)

	err = utils.DoWithTries(
		func() error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			dbpool, err = pgxpool.New(ctx, dsn)
			if err != nil {
				log.Panicf("Unable to create connection pool: %v\n", err)
				return err
			}
			// defer dbpool.Close()

			return nil
		}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Panic("Error with retires")
	}

	return dbpool, nil
}
