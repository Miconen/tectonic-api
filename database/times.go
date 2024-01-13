package database

import (
	"context"
	"strings"
	"tectonic-api/models"
	"time"
)

func SelectTime() (models.Time, error) {
	time := models.Time{}
	return time, nil
}

func InsertTime(f map[string]string) (models.Time, error) {
	teamIds := strings.Split(f["user_ids"], ",")

	timeData := map[string]interface{}{
		"boss_name": f["boss_name"],
		"date":      time.Now(),
		"time":      f["time"],
	}

	sql, args, err := psql.Insert("times").SetMap(timeData).Suffix("RETURNING *").ToSql()
	if err != nil {
		return models.Time{}, err
	}

	row := db.QueryRow(context.Background(), sql, args...)

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
	query := psql.Insert("teams").Columns("run_id", "guild_id", "user_id")

	for _, teammate := range res.Team {
		query = query.Values(teammate.RunId, teammate.GuildId, teammate.UserId)
	}

	sql, args, err = query.ToSql()
	_, err = db.Exec(context.Background(), sql, args...)
	if err != nil {
		return models.Time{}, err
	}

	return res, nil
}

func DeleteTime() {
}
