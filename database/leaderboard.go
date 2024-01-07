package database

import (
	"context"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectLeaderboard(filter map[string]string) (models.Users, error) {
	query := psql.Select("*").From("users").OrderBy("points DESC").Limit(50)

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
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
