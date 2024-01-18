package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func SelectUser(ctx context.Context, f map[string]string) (models.User, error) {
	query := psql.Select("users.*").From("users")

	switch {
	case f["rsn"] != "":
		query = query.Join("rsn ON users.user_id = rsn.user_id AND users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.rsn": f["rsn"]})
	case f["wom_id"] != "":
		query = query.Join("rsn ON users.user_id = rsn.user_id AND users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.wom_id": f["wom_id"]})
	case f["user_id"] != "":
		query = query.Where(squirrel.Eq{"users.user_id": f["user_id"], "users.guild_id": f["guild_id"]})
	default:
		return models.User{}, fmt.Errorf("no valid search key provided")
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.User{}, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, sql, args...)

	var user models.User

	err = row.Scan(&user.UserId, &user.GuildId, &user.Points)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

const ERROR_UNACTIVATED_GUILD string = "insert or update on table \"users\" violates foreign key constraint \"users_ibfk_1\""

func InsertUser(ctx context.Context, f map[string]string, wid string) error {
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

	query := psql.Insert("users").Columns("guild_id", "user_id").Values(f["guild_id"], f["user_id"])
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTagUser, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		pgxerr := err.(*pgconn.PgError)
		return fmt.Errorf("%s", pgxerr.Message)
	}

	query = psql.Insert("rsn").Columns("guild_id", "user_id", "rsn", "wom_id").Values(f["guild_id"], f["user_id"], f["rsn"], wid)
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

func DeleteUser(ctx context.Context, f map[string]string) error {
	query := psql.Delete("users")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}
