package database

import (
	"context"
	"fmt"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/Masterminds/squirrel"
)

func SelectGuild(f map[string]string) (models.Guild, error) {
	query := psql.Select("*").From("guilds")

	for key, value := range f {
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

func InsertGuild(g string) error {
	query := psql.Insert("guilds").Columns("guild_id").Values(g)
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

func DeleteGuild(f map[string]string) error {
	query := psql.Delete("guilds")

	for key, value := range f {
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

func UpdateGuild(g string, f map[string]string) error {
	// .SetMap requires an empty interface map so we convert here
	// As f passes parameters caught from HTTP requests, thus always string
	filter := utils.StringMapToInterfaceMap(f)

	query := psql.Update("guilds").SetMap(filter).Where(squirrel.Eq{"guild_id": g})
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
