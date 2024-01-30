package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectRsns(ctx context.Context, f map[string]string) ([]models.RSN, error) {
	query := psql.Select("*").From("rsn")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var rsns []models.RSN

	// Iterate through the rows and scan data into User struct
	for rows.Next() {
		var rsn models.RSN
		if err := rows.Scan(&rsn.UserId, &rsn.GuildId, &rsn.WomId, &rsn.RSN); err != nil {
			return nil, err
		}
		rsns = append(rsns, rsn)
	}

	return rsns, nil
}

func InsertRsn(ctx context.Context, f models.InputRSN, wid string) error {
	query := psql.Insert("rsn").Columns("guild_id", "user_id", "rsn", "wom_id").Values(f.GuildId, f.UserId, f.RSN, wid)
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

func DeleteRsn(ctx context.Context, f map[string]string) error {
	query := psql.Delete("rsn")

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
