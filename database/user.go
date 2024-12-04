package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/jackc/pgx/v5/pgconn"
)

func InsertUser(ctx context.Context, f models.InputUser, wid string) error {
	conn, err := pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Rollback the transaction if it hasn't been committed

	query := psql.Insert("users").Columns("guild_id", "user_id").Values(f.GuildId, f.UserId)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTagUser, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		pgxerr := err.(*pgconn.PgError)
		return fmt.Errorf("%s", pgxerr.Message)
	}

	query = psql.Insert("rsn").Columns("guild_id", "user_id", "rsn", "wom_id").Values(f.GuildId, f.UserId, f.RSN, wid)
	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}

	commandTagRsn, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	if commandTagUser.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTagUser.RowsAffected())
	} else if commandTagRsn.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTagRsn.RowsAffected())
	}

	return nil
}
