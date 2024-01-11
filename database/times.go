package database

import (
	"context"
	"tectonic-api/models"
	"time"
)

func SelectTime() (models.Time, error) {
	time := models.Time{}
	return time, nil
}

func InsertTime(ticks int64, boss string, teamData []map[string]interface{}) error {
	// Create time entry
	timeData := map[string]interface{}{
		"time":      ticks,
		"boss_name": boss,
		"date":      time.Now(),
	}

	sql, args, err := psql.Insert("times").SetMap(timeData).Suffix("RETURNING id").ToSql()
	if err != nil {
		return err
	}

	var time int
	err = db.QueryRow(context.Background(), sql, args...).Scan(&time)
	if err != nil {
		return err
	}

	// Update teamData to include the time_id
	for i := range teamData {
		teamData[i]["run_id"] = time
	}

	// Create team entries
	sql, args, err = psql.Insert("teams").Columns("user_id", "guild_id", "run_id").Values(teamData).ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql, args...)
	return err
}

func DeleteTime() {
}
