package database

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"strings"
=======
>>>>>>> main
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

<<<<<<< HEAD
type PointUpdate struct {
	Users       []models.User `json:"users"`
	GivenPoints int           `json:"given_points"`
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
=======
func UpdatePoints(ctx context.Context, f models.User) error {
	sql, args, err := psql.Update("users").Set("points", squirrel.Expr("points + ?", f.Points)).Where(squirrel.Eq{"user_id": f.UserId, "guild_id": f.GuildId}).ToSql()
	if err != nil {
		return err
>>>>>>> main
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

// Generate a function that returns updated users and Scan them into the []PointUpdate struct
func fetchUpdatedUsers(ctx context.Context, userIds []string, f map[string]string) ([]PointUpdate, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := psql.Select("users.user_id", "users.points", "point_sources.points").
		From("users").
		Where(squirrel.And{
			squirrel.Eq{"users.user_id": userIds},
			squirrel.Eq{"users.guild_id": guildID},
		}).
		Join("point_sources").
		Where(squirrel.And{
			squirrel.Eq{"point_sources.user_id": userIds},
			squirrel.Eq{"point_sources.guild_id": guildID},
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updates := PointUpdate{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Users.UserId, &u.Users.GuildId, &u.Users.Points)
		if err != nil {
			return nil, err
		}

		updates = append(updates, u)
	}

	return updates, nil
}
