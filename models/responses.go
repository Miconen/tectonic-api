package models

import (
	"encoding/json"
	"time"

	"tectonic-api/database"
)

// DetailedUser - returned by user lookup and WOM competition endpoints
type DetailedUser struct {
	UserId             string                  `json:"user_id"`
	GuildId            string                  `json:"guild_id"`
	Points             int                     `json:"points"`
	RSNs               []UserRsn               `json:"rsns"`
	Times              []UserTime              `json:"times"`
	Events             []UserEvent             `json:"events"`
	Achievements       []UserAchievement       `json:"achievements"`
	CombatAchievements []UserCombatAchievement `json:"combat_achievements"`
}

// User sub-models

type UserRsn struct {
	RSN   string `json:"rsn"`
	WomId string `json:"wom_id"`
}

func UserRsnsFromRows(rows []database.GetUserRsnsRow) []UserRsn {
	result := make([]UserRsn, len(rows))
	for i := range rows {
		result[i] = UserRsn{
			RSN:   rows[i].Rsn,
			WomId: rows[i].WomID,
		}
	}
	return result
}

type UserAchievement struct {
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	DiscordIcon string `json:"discord_icon"`
	Order       int16  `json:"order"`
}

func UserAchievementsFromRows(rows []database.GetUserAchievementsRow) []UserAchievement {
	result := make([]UserAchievement, len(rows))
	for i := range rows {
		result[i] = UserAchievement{
			Name:        rows[i].Name,
			Thumbnail:   rows[i].Thumbnail,
			DiscordIcon: rows[i].DiscordIcon,
		}
	}
	return result
}

type UserEvent struct {
	Name           string `json:"name"`
	WomID          string `json:"wom_id"`
	GuildID        string `json:"guild_id"`
	Placement      int16  `json:"placement"`
	PositionCutoff int16  `json:"position_cutoff"`
	Solo           bool   `json:"solo"`
}

func UserEventFromRows(rows []database.GetUserEventsRow) []UserEvent {
	result := make([]UserEvent, len(rows))
	for i := range rows {
		result[i] = UserEvent{
			Name:           rows[i].Name,
			WomID:          rows[i].EventID,
			GuildID:        rows[i].GuildID,
			Placement:      rows[i].Placement,
			PositionCutoff: rows[i].PositionCutoff,
			Solo:           rows[i].Solo,
		}
	}
	return result
}

type TimeTeammates struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

type UserTime struct {
	Id          int32           `json:"run_id"`
	BossName    string          `json:"boss_name"`
	DisplayName string          `json:"display_name"`
	Category    string          `json:"category"`
	Solo        bool            `json:"solo"`
	Date        time.Time       `json:"date"`
	Time        int32           `json:"time"`
	Teammates   []TimeTeammates `json:"team"`
}

func UserTimesFromRows(rows []database.GetUserTimesRow) []UserTime {
	if len(rows) == 0 {
		return []UserTime{}
	}

	result := make([]UserTime, 0)
	t := UserTime{
		Id:        0,
		Teammates: make([]TimeTeammates, 0),
	}

	for i := range rows {
		if rows[i].RunID != t.Id {
			if i != 0 {
				result = append(result, t)
			}
			t = UserTime{
				Id:          rows[i].RunID,
				BossName:    rows[i].BossName,
				DisplayName: rows[i].DisplayName,
				Category:    rows[i].Category,
				Solo:        rows[i].Solo,
				Date:        rows[i].Date.Time,
				Time:        rows[i].Time,
				Teammates:   make([]TimeTeammates, 0),
			}
		}
		t.Teammates = append(t.Teammates, TimeTeammates{
			UserID:  rows[i].UserID,
			GuildID: rows[i].GuildID,
		})
	}

	result = append(result, t)
	return result
}

type UserCombatAchievement struct {
	Name string `json:"name"`
}

func UserCombatAchievementsFromRows(rows []string) []UserCombatAchievement {
	result := make([]UserCombatAchievement, len(rows))
	for i := range rows {
		result[i] = UserCombatAchievement{
			Name: rows[i],
		}
	}
	return result
}

// Time creation response
type TimeResponse struct {
	BossName string `json:"boss_name"`
	Time     int    `json:"time"`
	OldTime  int    `json:"time_old"`
	RunID    int    `json:"run_id"`
}

// Event detail response
type DetailedEvent struct {
	Participations []EventParticipation `json:"participations"`
}

type EventParticipation struct {
	UserId    string `json:"user_id"`
	Placement int    `json:"placement"`
}

type GuildDetails struct {
	Teammates       []GuildTeammate      `json:"teammates,omitempty"`
	Pbs             []GuildPb            `json:"pbs,omitempty"`
	Bosses          []GuildBoss          `json:"bosses,omitempty"`
	Categories      []GuildCategory      `json:"categories,omitempty"`
	GuildBosses     []GuildBossEntry     `json:"guild_bosses,omitempty"`
	GuildCategories []GuildCategoryEntry `json:"guild_categories,omitempty"`
}

type Guild struct {
	GuildID      string  `json:"guild_id"`
	Multiplier   int32   `json:"multiplier"`
	PbChannelID  *string `json:"pb_channel_id"`
	ModChannelID *string `json:"mod_channel_id"`
	UserCount    int64   `json:"user_count"`
	TimeCount    int64   `json:"time_count"`

	GuildDetails
}

type GuildTeammate struct {
	RunID   int32  `json:"run_id"`
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

type GuildPb struct {
	RunID    int32  `json:"run_id"`
	Time     int32  `json:"time"`
	BossName string `json:"boss_name"`
	Date     string `json:"date"`
	GuildID  string `json:"guild_id"`
}

type GuildBoss struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	Solo        bool   `json:"solo"`
}

type GuildCategory struct {
	Thumbnail string `json:"thumbnail"`
	Order     int16  `json:"order"`
	Name      string `json:"name"`
}

type GuildBossEntry struct {
	Boss     string `json:"boss"`
	GuildID  string `json:"guild_id"`
	PbID     *int32 `json:"pb_id"`
	Category string `json:"category"`
}

type GuildCategoryEntry struct {
	GuildID   string `json:"guild_id"`
	Category  string `json:"category"`
	MessageID string `json:"message_id"`
}

func GuildResponseFromRow(row database.GetGuildRow) Guild {
	var pbChannelID *string
	if row.PbChannelID.Valid {
		pbChannelID = &row.PbChannelID.String
	}

	var modChannelID *string
	if row.ModChannelID.Valid {
		modChannelID = &row.ModChannelID.String
	}

	return Guild{
		GuildID:      row.GuildID,
		Multiplier:   row.Multiplier,
		PbChannelID:  pbChannelID,
		ModChannelID: modChannelID,
		UserCount:    row.UserCount,
		TimeCount:    row.TimeCount,
	}
}

func GuildResponseFromDetailedRow(row database.GetDetailedGuildRow) Guild {
	var pbChannelID *string
	if row.PbChannelID.Valid {
		pbChannelID = &row.PbChannelID.String
	}

	var modChannelID *string
	if row.ModChannelID.Valid {
		modChannelID = &row.ModChannelID.String
	}

	g := Guild{
		GuildID:      row.GuildID,
		Multiplier:   row.Multiplier,
		PbChannelID:  pbChannelID,
		ModChannelID: modChannelID,
		UserCount:    row.UserCount,
		TimeCount:    row.TimeCount,
		GuildDetails: GuildDetails{
			Teammates:       []GuildTeammate{},
			Pbs:             []GuildPb{},
			Bosses:          []GuildBoss{},
			Categories:      []GuildCategory{},
			GuildBosses:     []GuildBossEntry{},
			GuildCategories: []GuildCategoryEntry{},
		},
	}

	json.Unmarshal(row.Teammates, &g.Teammates)
	json.Unmarshal(row.Pbs, &g.Pbs)
	json.Unmarshal(row.Bosses, &g.Bosses)
	json.Unmarshal(row.Categories, &g.Categories)
	json.Unmarshal(row.GuildBosses, &g.GuildBosses)
	json.Unmarshal(row.GuildCategories, &g.GuildCategories)

	return g
}

type LeaderboardUser struct {
	UserID  string    `json:"user_id"`
	GuildID string    `json:"guild_id"`
	Points  int32     `json:"points"`
	RSNs    []UserRsn `json:"rsns"`
}

func LeaderboardFromRows(rows []database.GetLeaderboardRow) []LeaderboardUser {
	list := make([]LeaderboardUser, 0, len(rows))
	for _, row := range rows {
		user := LeaderboardUser{
			UserID:  row.UserID,
			GuildID: row.GuildID,
			Points:  row.Points,
			RSNs:    []UserRsn{},
		}
		json.Unmarshal(row.Rsns, &user.RSNs)
		list = append(list, user)
	}
	return list
}
