package database

import (
	"context"
	"strings"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectUsers(f map[string]string) (models.Users, error) {
	userIds := strings.Split(f["user_ids"], ",")

	query := squirrel.Select("*").From("users").Where(squirrel.Eq{"guild_id": f["guild_id"]})

	for _, value := range userIds {
		query = query.Where(squirrel.Eq{"user_id": value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Users{}, err
	}

	// Executing the query
	rows, err := db.Query(context.Background(), sql, args...)
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
