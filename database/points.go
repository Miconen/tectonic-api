package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func UpdatePoints(ctx context.Context, f map[string]string) error {
	sql, args, err := psql.Update("users").Set("points", squirrel.Expr("points + ?", f["points"])).Where(squirrel.Eq{"user_id": f["user_id"], "guild_id": f["guild_id"]}).ToSql()
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
