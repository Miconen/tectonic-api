package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func UpdatePoints(ctx context.Context, f models.User) error {
	sql, args, err := psql.Update("users").Set("points", squirrel.Expr("points + ?", f.Points)).Where(squirrel.Eq{"user_id": f.UserId, "guild_id": f.GuildId}).ToSql()
	if err != nil {
		return err
	}

	fmt.Println(sql, args)

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
