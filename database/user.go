package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectUser(filter map[string]string) (models.User, error) {
	query := psql.Select("*").From("users")

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.User{}, err
	}

	row := db.QueryRow(context.Background(), sql, args...)

	var user models.User

	err = row.Scan(&user.UserId, &user.GuildId, &user.Points)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func InsertUser(user models.User) error {
	query := psql.Insert("users").Columns("guild_id", "user_id").Values(user.GuildId, user.UserId)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}

func DeleteUser(filter map[string]string) error {
	query := psql.Delete("users")

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}
