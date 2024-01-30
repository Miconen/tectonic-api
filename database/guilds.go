package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

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

func InsertGuild(ctx context.Context, g models.InputGuild) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := psql.Insert("guilds").Columns("guild_id", "multiplier", "pb_channel_id").Values(g.GuildId, g.Multiplier, g.PbChannelId)
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

	sql = `INSERT INTO guild_categories(guild_id, category, message_id)
		   SELECT $1, name, '' FROM categories`
	_, err = tx.Exec(ctx, sql, g.GuildId)
	if err != nil {
		return err
	}

	sql = `INSERT INTO guild_bosses(guild_id, boss, pb_id)
		   SELECT $1, name, NULL FROM bosses`
	_, err = tx.Exec(ctx, sql, g.GuildId)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	if rows := commandTagGuild.RowsAffected(); rows != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows)
	}

	return nil
}

func DeleteGuild(ctx context.Context, f map[string]string) error {
	query := psql.Delete("guilds").Where(squirrel.Eq{"guild_id": f["guild_id"]})

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

func UpdateGuild(ctx context.Context, g models.UpdateGuild) error {
	query := psql.Update("guilds").Set("multiplier", g.Multiplier).Set("pb_channel_id", g.PbChannelId).Where(squirrel.Eq{"guild_id": g.GuildId})
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
