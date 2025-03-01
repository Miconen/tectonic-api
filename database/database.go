package database

import (
	"context"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ERROR_UNACTIVATED_GUILD string = "insert or update on table \"users\" violates foreign key constraint \"users_ibfk_1\""

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
var pool *pgxpool.Pool

func InitDB() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	pool = conn // Store the connection in a package-level variable

	return conn, nil
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
