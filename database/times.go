package database

import (
	"context"
	"fmt"
	"strings"
	"tectonic-api/models"
	"time"

	"github.com/Masterminds/squirrel"
)

func SelectTime() (models.Time, error) {
	time := models.Time{}
	return time, nil
}

func CheckPb(ctx context.Context, f map[string]string) (int, error) {
	query := psql.Select("t.time").
		From("guild_bosses gb").
		Join("times t ON gb.pb_id = t.run_id").
		Where(squirrel.Eq{"gb.guild_id": f["guild_id"], "gb.boss": f["boss"]})
	sql, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, sql, args...)

	var pb int

	err = row.Scan(&pb)
	if err != nil {
		return 0, err
	}

	return pb, nil
}

func InsertTime(ctx context.Context, f map[string]string) (models.Time, error) {
	teamIds := strings.Split(f["user_ids"], ",")

	timeData := map[string]interface{}{
		"boss_name": f["boss"],
		"date":      time.Now(),
		"time":      f["time"],
	}

	sql, args, err := psql.Insert("times").SetMap(timeData).Suffix("RETURNING *").ToSql()
	if err != nil {
		return models.Time{}, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return models.Time{}, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, sql, args...)

	var runId int
	res := models.Time{Team: []models.Teammate{}}
	err = row.Scan(&res.Time, &res.BossName, &runId, &res.Date)
	if err != nil {
		return models.Time{}, err
	}

	res.RunId = runId

	// Update teamData to include the time_id
	for _, v := range teamIds {
		res.Team = append(res.Team, models.Teammate{
			RunId:   runId,
			GuildId: f["guild_id"],
			UserId:  v,
		})
	}

	// Update guild pb (guild_bosses.pb_id)
	sql, args, err = psql.Update("guild_bosses").
		Set("pb_id", runId).
		Where(squirrel.Eq{"guild_id": f["guild_id"], "boss": f["boss"]}).
		ToSql()
	if err != nil {
		return models.Time{}, err
	}

	commandTag, err := conn.Exec(ctx, sql, args...)
	if err != nil {
		return models.Time{}, err
	}

	if commandTag.RowsAffected() != 1 {
		return models.Time{}, fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	query := psql.Insert("teams").Columns("run_id", "guild_id", "user_id")

	for _, teammate := range res.Team {
		query = query.Values(teammate.RunId, teammate.GuildId, teammate.UserId)
	}

	sql, args, err = query.ToSql()
	if err != nil {
		return models.Time{}, err
	}

	_, err = conn.Exec(ctx, sql, args...)
	if err != nil {
		return models.Time{}, err
	}

	return res, nil
}

func DeleteTime(ctx context.Context, f map[string]string) error {
	sql, args, err := psql.Delete("users").Where(squirrel.Eq{"run_id": f["time_id"]}).ToSql()
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
