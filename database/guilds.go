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
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background()) // Rollback the transaction if it hasn't been committed

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

	err = InitializeGuildCategories(g)
	if err != nil {
		return err
	}

	err = InitializeGuildBosses(g)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
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

	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
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

func InitializeGuildBosses(g string) error {
	sql, args, err := psql.Select("name").From("bosses").ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	bosses := make([]map[string]interface{}, 0)

	for rows.Next() {
		var b string
		// Assuming the columns are named "column1", "column2", etc.
		if err := rows.Scan(&b); err != nil {
			return err
		}

		boss := map[string]interface{}{
			"boss":     b,
			"guild_id": g,
			"pb_id":    nil,
		}

		bosses = append(bosses, boss)
	}

	// Insert data into guild_bosses
	query := psql.Insert("guild_bosses").Columns("boss", "guild_id", "pb_id")

	for _, v := range bosses {
		query = query.Values(v["boss"], v["guild_id"], v["pb_id"])
	}

	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func InitializeGuildCategories(g string) error {
	sql, args, err := psql.Select("name").From("categories").ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	categories := make([]map[string]interface{}, 0)

	for rows.Next() {
		var c string
		// Assuming the columns are named "column1", "column2", etc.
		if err := rows.Scan(&c); err != nil {
			return err
		}

		category := map[string]interface{}{
			"guild_id":   g,
			"category":   c,
			"message_id": "",
		}

		categories = append(categories, category)
	}

	// Insert data into guild_bosses
	query := psql.Insert("guild_categories").Columns("guild_id", "category", "message_id")

	for _, v := range categories {
		query = query.Values(v["guild_id"], v["category"], v["message_id"])
	}

	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}
