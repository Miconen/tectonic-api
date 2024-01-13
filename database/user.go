package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectUser(f map[string]string) (models.User, error) {
	query := psql.Select("users.*").From("users")

	switch {
	case f["rsn"] != "":
		query = query.Join("rsn ON users.user_id = rsn.user_id AND users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.rsn": f["rsn"]})
	case f["wom_id"] != "":
		query = query.Join("rsn ON users.user_id = rsn.user_id AND users.guild_id = rsn.guild_id").Where(squirrel.Eq{"rsn.wom_id": f["wom_id"]})
	case f["user_id"] != "":
		query = query.Where(squirrel.Eq{"users.user_id": f["user_id"], "users.guild_id": f["guild_id"]})
	default:
		return models.User{}, fmt.Errorf("No valid search key provided")
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

func InsertUser(f map[string]string, wid string) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // Rollback the transaction if it hasn't been committed

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

	err = InsertRsn(f, wid)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
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
