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

type DetailedGuildJSON struct {
	GuildID         string          `json:"guild_id"`
	PbChannelID     pgtype.Text     `json:"pb_channel_id"`
	Teammates       json.RawMessage `json:"teammates"`
	Pbs             json.RawMessage `json:"pbs"`
	Bosses          json.RawMessage `json:"bosses"`
	Categories      json.RawMessage `json:"categories"`
	GuildBosses     json.RawMessage `json:"guild_bosses"`
	GuildCategories json.RawMessage `json:"guild_categories"`
}

func NewDetailedGuildFromRow(row GetDetailedGuildRow) DetailedGuildJSON {
	return DetailedGuildJSON{
		GuildID:         row.GuildID,
		PbChannelID:     row.PbChannelID,
		Teammates:       row.Teammates,
		Pbs:             row.Pbs,
		Bosses:          row.Bosses,
		Categories:      row.Categories,
		GuildBosses:     row.GuildBosses,
		GuildCategories: row.GuildCategories,
	}
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
