package database

import (
	"context"
	"strings"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectUsers(ctx context.Context, g string, f map[string]string) (models.Users, error) {
	query := psql.Select("*").From("users").Where(squirrel.Eq{"guild_id": g})

	if f["user_ids"] != "" {
		userIds := strings.Split(f["user_ids"], ",")
		query = query.Where(squirrel.Eq{"user_id": userIds})
	}
	if f["wom_ids"] != "" {
		womIds := strings.Split(f["wom_ids"], ",")
		query = query.Join("rsn ON users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.wom_id": womIds})
	}
	if f["rsns"] != "" {
		rsns := strings.Split(f["rsns"], ",")
		query = query.Join("rsn ON users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.rsn": rsns})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Users{}, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return models.Users{}, err
	}
	defer conn.Release()

	// Executing the query
	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return models.Users{}, err
	}
	defer rows.Close()

	var users models.Users

	// Iterate through the rows and scan data into User struct
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.GuildId, &user.Points); err != nil {
			return models.Users{}, err
		}
		users.Users = append(users.Users, user)
	}

	if err := rows.Err(); err != nil {
		return models.Users{}, err
	}

	return users, nil
}
