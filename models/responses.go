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
	Rank               int64                   `json:"rank"`
	Tier               *UserTier               `json:"tier,omitempty"`
	RSNs               []UserRsn               `json:"rsns"`
	Records            []UserRecord            `json:"records"`
	Events             []UserEvent             `json:"events"`
	Achievements       []UserAchievement       `json:"achievements"`
	CombatAchievements []UserCombatAchievement `json:"combat_achievements"`
}

// UserTier - the user's current rank tier based on points
type UserTier struct {
	Name         string  `json:"name"`
	Icon         *string `json:"icon,omitempty"`
	RoleID       *string `json:"role_id,omitempty"`
	MinPoints    int32   `json:"min_points"`
	DisplayOrder int16   `json:"display_order"`
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

type RecordTeammate struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

type UserRecord struct {
	Id          int32            `json:"record_id"`
	BossName    string           `json:"boss_name"`
	DisplayName string           `json:"display_name"`
	Category    string           `json:"category"`
	Solo        bool             `json:"solo"`
	ValueType   string           `json:"value_type"`
	Date        time.Time        `json:"date"`
	Value       int32            `json:"value"`
	Teammates   []RecordTeammate `json:"team"`
}

func UserRecordsFromRows(rows []database.GetUserRecordsRow) []UserRecord {
	if len(rows) == 0 {
		return []UserRecord{}
	}

	result := make([]UserRecord, 0)
	r := UserRecord{
		Id:        0,
		Teammates: make([]RecordTeammate, 0),
	}

	for i := range rows {
		if rows[i].RecordID != r.Id {
			if i != 0 {
				result = append(result, r)
			}
			r = UserRecord{
				Id:          rows[i].RecordID,
				BossName:    rows[i].BossName,
				DisplayName: rows[i].DisplayName,
				Category:    rows[i].Category,
				Solo:        rows[i].Solo,
				ValueType:   rows[i].ValueType,
				Date:        rows[i].Date.Time,
				Value:       rows[i].Value,
				Teammates:   make([]RecordTeammate, 0),
			}
		}
		r.Teammates = append(r.Teammates, RecordTeammate{
			UserID:  rows[i].UserID,
			GuildID: rows[i].GuildID,
		})
	}

	result = append(result, r)
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

// Record creation response
type RecordResponse struct {
	BossName string `json:"boss_name"`
	Value    int    `json:"value"`
	OldValue int    `json:"value_old"`
	RecordID int    `json:"record_id"`
	Position *int   `json:"position,omitempty"`
}

// Event detail response
type DetailedEvent struct {
	Participations []EventParticipation `json:"participations"`
}

type EventParticipation struct {
	UserId    string `json:"user_id"`
	Placement int    `json:"placement"`
}

// Guild response models

type GuildTeammate struct {
	RecordID int32  `json:"record_id"`
	UserID   string `json:"user_id"`
	GuildID  string `json:"guild_id"`
}

type GuildRecord struct {
	RecordID int32  `json:"record_id"`
	Value    int32  `json:"value"`
	BossName string `json:"boss_name"`
	Date     string `json:"date"`
	GuildID  string `json:"guild_id"`
	Position int64  `json:"position"`
}

type GuildBoss struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	Solo        bool   `json:"solo"`
	ValueType   string `json:"value_type"`
}

type GuildCategory struct {
	Thumbnail string `json:"thumbnail"`
	Order     int16  `json:"order"`
	Name      string `json:"name"`
}

type GuildBossEntry struct {
	Boss     string `json:"boss"`
	GuildID  string `json:"guild_id"`
	Category string `json:"category"`
}

type GuildCategoryEntry struct {
	GuildID   string `json:"guild_id"`
	Category  string `json:"category"`
	MessageID string `json:"message_id"`
}

type GuildRankResponse struct {
	Name         string  `json:"name"`
	MinPoints    int32   `json:"min_points"`
	Icon         *string `json:"icon,omitempty"`
	RoleID       *string `json:"role_id,omitempty"`
	DisplayOrder int16   `json:"display_order"`
}

type GuildDetails struct {
	Teammates       []GuildTeammate      `json:"teammates,omitempty"`
	Records         []GuildRecord        `json:"records,omitempty"`
	Bosses          []GuildBoss          `json:"bosses,omitempty"`
	Categories      []GuildCategory      `json:"categories,omitempty"`
	GuildBosses     []GuildBossEntry     `json:"guild_bosses,omitempty"`
	GuildCategories []GuildCategoryEntry `json:"guild_categories,omitempty"`
}

type GuildResponse struct {
	GuildID       string  `json:"guild_id"`
	Multiplier    int32   `json:"multiplier"`
	PbChannelID   *string `json:"pb_channel_id"`
	ModChannelID  *string `json:"mod_channel_id"`
	PositionCount int16   `json:"position_count"`
	UserCount     int64   `json:"user_count"`
	RecordCount   int64   `json:"record_count"`

	GuildDetails
}

func GuildResponseFromRow(row database.GetGuildRow) GuildResponse {
	var pbChannelID *string
	if row.PbChannelID.Valid {
		pbChannelID = &row.PbChannelID.String
	}

	var modChannelID *string
	if row.ModChannelID.Valid {
		modChannelID = &row.ModChannelID.String
	}

	return GuildResponse{
		GuildID:       row.GuildID,
		Multiplier:    row.Multiplier,
		PbChannelID:   pbChannelID,
		ModChannelID:  modChannelID,
		PositionCount: row.PositionCount,
		UserCount:     row.UserCount,
		RecordCount:   row.RecordCount,
	}
}

func GuildResponseFromDetailedRow(row database.GetDetailedGuildRow) GuildResponse {
	var pbChannelID *string
	if row.PbChannelID.Valid {
		pbChannelID = &row.PbChannelID.String
	}

	var modChannelID *string
	if row.ModChannelID.Valid {
		modChannelID = &row.ModChannelID.String
	}

	g := GuildResponse{
		GuildID:       row.GuildID,
		Multiplier:    row.Multiplier,
		PbChannelID:   pbChannelID,
		ModChannelID:  modChannelID,
		PositionCount: row.PositionCount,
		UserCount:     row.UserCount,
		RecordCount:   row.RecordCount,
		GuildDetails: GuildDetails{
			Teammates:       []GuildTeammate{},
			Records:         []GuildRecord{},
			Bosses:          []GuildBoss{},
			Categories:      []GuildCategory{},
			GuildBosses:     []GuildBossEntry{},
			GuildCategories: []GuildCategoryEntry{},
		},
	}

	json.Unmarshal(row.Teammates, &g.Teammates)
	json.Unmarshal(row.Records, &g.Records)
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

func GuildRanksFromRows(rows []database.GetGuildRanksRow) []GuildRankResponse {
	result := make([]GuildRankResponse, len(rows))
	for i, row := range rows {
		r := GuildRankResponse{
			Name:         row.Name,
			MinPoints:    row.MinPoints,
			DisplayOrder: row.DisplayOrder,
		}
		if row.Icon.Valid {
			r.Icon = &row.Icon.String
		}
		if row.RoleID.Valid {
			r.RoleID = &row.RoleID.String
		}
		result[i] = r
	}
	return result
}
