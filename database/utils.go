package database

import "github.com/jackc/pgx/v5/pgtype"

type GameUser struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	RSNs    []Rsn  `json:"rsn"`
}

type DetailedTime struct {
	Time      int32            `json:"time"`
	BossName  string           `json:"boss_name"`
	RunID     int32            `json:"run_id"`
	Date      pgtype.Timestamp `json:"date"`
	Teammates []GameUser       `json:"team"`
}

type DetailedUser struct {
	User  GameUser `json:"user"`
	Times []Time   `json:"time"`
}
