package pgdb

import (
	"art_space/internal/envvar"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PGDB interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Begin(context.Context) (pgx.Tx, error)
}

type Queries struct {
	db PGDB
}

func New(db PGDB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{db: tx}
}

func NewDB(ctx context.Context) *pgx.Conn {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		envvar.Config.DB.Username, envvar.Config.DB.Password,
		envvar.Config.DB.Host, envvar.Config.DB.Port, envvar.Config.DB.Name,
	)

	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Не удалось пингануть БД: %v", err)
	}

	return db
}
