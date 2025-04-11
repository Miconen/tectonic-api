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

type ConstraintType uint

const (
	PrimaryKey ConstraintType = iota
	ForeignKey
	Unique
)

type ConstraintDetail struct {
	Type         ConstraintType
	Table        string
	ForeignTable string
}

func GetConstraintsTable(ctx context.Context, conn *pgxpool.Conn) (map[string]ConstraintDetail, error) {
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}

	result := make(map[string]ConstraintDetail)
	rows, err := pool.Query(ctx, `
		SELECT 
			tc.constraint_name,
			tc.table_name,
			ccu.table_name AS foreign_table_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.table_schema = 'public'
			AND tc.constraint_name NOT LIKE '%not_null'
			AND tc.constraint_name NOT LIKE 'goose%'
		`,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var key, table, foreign string
		if err = rows.Err(); err != nil {
			return nil, err
		}

		if err = rows.Scan(&key, &table, &foreign); err != nil {
			return nil, err
		}

		var ct ConstraintType
		if table == foreign {
			ct = PrimaryKey
		} else {
			ct = ForeignKey
		}

		result[key] = ConstraintDetail{
			Type:         ct,
			Table:        table,
			ForeignTable: foreign,
		}
	}

	return result, nil
}

func AcquireConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return pool.Acquire(ctx)
}

func CreateTx(ctx context.Context) (pgx.Tx, error) {
	return pool.BeginTx(ctx, pgx.TxOptions{})
}
