package database

import (
	"context"
	"fmt"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/Masterminds/squirrel"
)

func SelectGuild(ctx context.Context, f map[string]string) (models.Guild, error) {
	query := psql.Select("*").From("guilds")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Guild{}, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return models.Guild{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, sql, args...)

	var guild models.Guild

	err = row.Scan(&guild.GuildId, &guild.Multiplier, &guild.PbChannelId)
	if err != nil {
		return models.Guild{}, err
	}

	return guild, nil
}

func InsertGuild(ctx context.Context, g string) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := psql.Insert("guilds").Columns("guild_id").Values(g)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTagGuild, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	query = psql.Insert("guild_categories").
		Columns("guild_id", "category", "message_id").
		Suffix("SELECT $1, name, '' FROM categories", g)

	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}

	commandTagCategories, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	query = psql.Insert("guild_bosses").
		Columns("guild_id", "boss", "pb_id").
		Suffix("SELECT $1, name, NULL FROM bosses", g)

	sql, args, err = query.ToSql()
	if err != nil {
		return err
	}

	commandTagBosses, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	if rows := commandTagGuild.RowsAffected(); rows != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows)
	} else if rows := commandTagCategories.RowsAffected(); rows != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows)
	} else if rows := commandTagBosses.RowsAffected(); rows != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows)
	}

	return nil
}

func DeleteGuild(ctx context.Context, f map[string]string) error {
	query := psql.Delete("guilds")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}

func UpdateGuild(ctx context.Context, g string, f map[string]string) error {
	// .SetMap requires an empty interface map so we convert here
	// As f passes parameters caught from HTTP requests, thus always string
	filter := utils.StringMapToInterfaceMap(f)

	query := psql.Update("guilds").SetMap(filter).Where(squirrel.Eq{"guild_id": g})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}
