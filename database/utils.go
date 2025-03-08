package database

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type DetailedRsn struct {
	Rsn   string `json:"rsn"`
	WomID string `json:"wom_id"`
}

type UserData struct {
	UserID  string          `json:"user_id"`
	GuildID string          `json:"guild_id"`
	Points  int32           `json:"points"`
	RSNs    json.RawMessage `json:"rsns"`
}

type DetailedTime struct {
	Time      int32            `json:"time"`
	BossName  string           `json:"boss_name"`
	RunID     int32            `json:"run_id"`
	Date      pgtype.Timestamp `json:"date"`
	Teammates []UserData       `json:"team"`
}

// TODO: Probably separate this in multiple queries, my head hurts while
// making sql queries...
type DetailedUser struct {
	UserID  string          `json:"user_id"`
	GuildID string          `json:"guild_id"`
	Points  int32           `json:"points"`
	RSNs    json.RawMessage `json:"rsns"`
	Times   json.RawMessage `json:"times"`
}

func NewDetailedUserFromRows(rows []GetDetailedUsersRow) []DetailedUser {
	list := make([]DetailedUser, 0, len(rows))

	for _, row := range rows {
		user := DetailedUser{
			UserID:  row.UserID,
			GuildID: row.GuildID,
			Points:  row.Points,
			RSNs:    row.Rsns,
			Times:   row.Times,
		}

		list = append(list, user)
	}

	return list
}

func NewLeaderboardFromRows(rows []GetLeaderboardRow) []UserData {
	list := make([]UserData, 0, len(rows))

	for _, row := range rows {
		user := UserData{
			UserID:  row.UserID,
			GuildID: row.GuildID,
			Points:  row.Points,
			RSNs:    row.Rsns,
		}

		list = append(list, user)
	}

	return list
}
