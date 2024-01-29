package database

import (
	"context"
	"fmt"
	"strings"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

type PointUpdate struct {
	User        models.User `json:"user"`
	OldPoints   int         `json:"old_points"`
	GivenPoints int         `json:"given_points"`
}

func getUpdateSubquery(f map[string]string, pkey string) squirrel.Sqlizer {
	if pkey == "point_event" {
		// Sub query to select the value to increment user points with
		return squirrel.Expr("points + (?)", squirrel.Select("points").
			From("point_sources").
			Where(squirrel.And{
				squirrel.Eq{"source": f["point_event"]},
				squirrel.Eq{"guild_id": f["guild_id"]},
			}),
		)
	}

	return squirrel.Expr("points + (?)", f["points"])
}

func UpdatePoints(ctx context.Context, f map[string]string, pkey string) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	userIds := strings.Split(f["user_ids"], ",")

	subquery := getUpdateSubquery(f, pkey)

	query := psql.Update("users").
		Set("points", subquery).
		Where(squirrel.And{
			squirrel.Eq{"user_id": userIds},
			squirrel.Eq{"guild_id": f["guild_id"]},
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() <= 0 {
		return fmt.Errorf("expected rows to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}
