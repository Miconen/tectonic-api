package database

import (
	"context"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func AcquireConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return pool.Acquire(ctx)
}

func CreateTx(ctx context.Context) (pgx.Tx, error) {
	return pool.BeginTx(ctx, pgx.TxOptions{})
}
