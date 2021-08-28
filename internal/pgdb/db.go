package pgdb

import (
	"art_space/internal/envvar"
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
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

func NewDB(ctx context.Context, conf *envvar.VaultConfiguration) *pgx.Conn {
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			logrus.Fatalf("Не удалось получить значение конфигурации для %s: %v", v, err)
		}

		return res
	}

	dbUsername := get("PGDB_USERNAME")
	dbUserPassword := get("PGDB_PASSWORD")
	dbHost := get("PGDB_HOST")
	dbPort := get("PGDB_PORT")
	dbName := get("PGDB_NAME")
	// dbSslMode := get("PGDB_SSLMODE")

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUsername, dbUserPassword, dbHost, dbPort, dbName,
	)

	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		logrus.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		logrus.Fatalf("Не удалось пингануть БД: %v", err)
	}

	// if err := initPostRequests(ctx, db); err != nil {
	// 	logrus.Fatal(err)
	// }

	// if err := initCommentRequests(ctx, db); err != nil {
	// 	logrus.Fatal(err)
	// }

	return db
}
