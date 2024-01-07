package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectGuild(filter map[string]string) (models.Guild, error) {
	query := psql.Select("*").From("guilds")

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Guild{}, err
	}

	row := db.QueryRow(context.Background(), sql, args...)

	var guild models.Guild

	err = row.Scan(&guild.GuildId, &guild.Multiplier, &guild.PbChannelId)
	if err != nil {
		return models.Guild{}, err
	}

	return guild, nil
}

func InsertGuild(guild models.Guild) error {
	query := psql.Insert("guilds").Columns("guild_id").Values(guild.GuildId)
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

func DeleteGuild(filter map[string]string) error {
	query := psql.Delete("guilds")

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

func UpdateGuild(g string, c map[string]interface{}) error {
	query := psql.Update("guilds").SetMap(c).Where(squirrel.Eq{"guild_id": g})
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
