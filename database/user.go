package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectUser(f map[string]string) (models.User, error) {
	query := psql.Select("*").From("users")

	for key, value := range f {
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

func InsertUser(f map[string]string) error {
	query := psql.Insert("users").Columns("guild_id", "user_id").Values(f["guild_id"], f["user_id"])
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

func DeleteUser(f map[string]string) error {
	query := psql.Delete("users")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

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
