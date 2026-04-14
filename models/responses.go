package models

import (
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
