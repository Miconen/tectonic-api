package models

import (
	"tectonic-api/database"
	"time"
)

// InputGuild Model - for creating guilds
// @Description Model of new guild data
type InputGuild struct {
	GuildId     string `json:"guild_id" validate:"required,discord_snowflake"`
	Multiplier  int    `json:"multiplier" validate:"omitempty,min=1,max=10"`
	PbChannelId string `json:"pb_channel_id" validate:"omitempty,discord_snowflake"`
}

// UpdateGuild Model - for updating guilds
// @Description Model of updated guild data
type UpdateGuild struct {
	GuildId     string `json:"guild_id" validate:"required,discord_snowflake"`
	Multiplier  int    `json:"multiplier" validate:"omitempty,min=1,max=10"`
	PbChannelId string `json:"pb_channel_id" validate:"omitempty,discord_snowflake"`
}

// CreateUserBody Model - for creating users
// @Description Model of new active guild member
type CreateUserBody struct {
	UserId string `json:"user_id" validate:"required,discord_snowflake"`
	RSN    string `json:"rsn" validate:"required,rsn,min=1,max=12"`
}

// CreateRsnBody Model - for adding RSNs to existing users
// @Description Model of new guild member RSN
type CreateRsnBody struct {
	RSN string `json:"rsn" validate:"required,rsn,min=1,max=12"`
}

// InputUser Model - legacy model (if still used)
// @Description Model of new active guild member
type InputUser struct {
	UserId  string `json:"user_id" validate:"required,discord_snowflake"`
	GuildId string `json:"guild_id" validate:"required,discord_snowflake"`
	RSN     string `json:"rsn" validate:"required,rsn,min=1,max=12"`
}

// InputTime Model - for creating time records
// @Description Model of a new time submission
type InputTime struct {
	Time     int      `json:"time" validate:"required,positive_time"`
	BossName string   `json:"boss_name" validate:"required,min=1,max=50"`
	UserIds  []string `json:"user_ids" validate:"required,min=1,max=8,dive,discord_snowflake"`
}

// InputUser Model - for creating user team records
// @Description Model of a simple user team record
type InputTeammate struct {
	UserId  string `json:"user_id" validate:"required,discord_snowflake"`
	GuildId string `json:"guild_id" validate:"required,discord_snowflake"`
}

// InputRSN Model - legacy RSN input (if still used)
// @Description Model of RSN association
type InputRSN struct {
	RSN     string `json:"rsn" validate:"required,rsn,min=1,max=12"`
	UserId  string `json:"user_id" validate:"required,discord_snowflake"`
	GuildId string `json:"guild_id" validate:"required,discord_snowflake"`
}

// InputEvent Model - for creating/registering events
// @Description Model of event registration data
type InputEvent struct {
	EventId        int      `json:"event_id" validate:"required,min=1"`
	TeamNames      []string `json:"team_names" validate:"omitempty,dive,min=1,max=64"`
	PositionCutoff int      `json:"position_cutoff" validate:"omitempty,min=1,max=3"`
}

// CategoryMessage Model - for guild category message updates
// @Description Model for category message ID updates
type CategoryMessage struct {
	MessageID string `json:"message_id" validate:"required,discord_snowflake"`
	Category  string `json:"category" validate:"required,min=1,max=30"`
}

// GuildParams Model - for complex guild updates (from guild_handler.go)
// @Description Model for guild parameter updates
type GuildParams struct {
	Multiplier       *float64          `json:"multiplier" validate:"omitempty,min=0.1,max=10"`
	PbChannelID      string            `json:"pb_channel_id" validate:"omitempty,discord_snowflake"`
	CategoryMessages []CategoryMessage `json:"category_messages" validate:"omitempty,dive"`
}

// Additional models that might need validation:

// PointUpdate Model - for manual point updates
// @Description Model for manual point adjustments
type PointUpdate struct {
	Points  int32    `json:"points" validate:"required"`
	UserIds []string `json:"user_ids" validate:"required,min=1,dive,discord_snowflake"`
	Reason  string   `json:"reason" validate:"omitempty,max=100"`
}

// CompetitionRequest Model - for WOM competition handling
// @Description Model for competition processing requests
type CompetitionRequest struct {
	CompetitionId int `json:"competition_id" validate:"required,min=1"`
	Cutoff        int `json:"cutoff" validate:"required,min=0"`
}

// AchievementInput Model - for achievement management
// @Description Model for achievement creation/updates
type AchievementInput struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Thumbnail   string `json:"thumbnail" validate:"required,url"`
	DiscordIcon string `json:"discord_icon" validate:"required,min=1,max=100"`
	Order       int16  `json:"order" validate:"min=0"`
}

// BulkUserOperation Model - for operations on multiple users
// @Description Model for bulk user operations
type BulkUserOperation struct {
	UserIds   []string               `json:"user_ids" validate:"required,min=1,max=50,dive,discord_snowflake"`
	Operation string                 `json:"operation" validate:"required,oneof=add remove update"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// TimeUpdate Model - for updating existing times
// @Description Model for time record updates
type TimeUpdate struct {
	Time     *int     `json:"time" validate:"omitempty,positive_time"`
	BossName *string  `json:"boss_name" validate:"omitempty,min=1,max=50"`
	UserIds  []string `json:"user_ids" validate:"omitempty,min=1,max=8,dive,discord_snowflake"`
}

// User Model
// @Description Model of active guild member
type User struct {
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
	Points  int    `json:"points"`
}

// Detailed User Model
// @Description Model of active guild member containing events, times and achievements
type DetailedUser struct {
	UserId       string            `json:"user_id"`
	GuildId      string            `json:"guild_id"`
	Points       int               `json:"points"`
	RSNs         []UserRsn         `json:"rsns"`
	Times        []UserTime        `json:"times"`
	Events       []UserEvent       `json:"events"`
	Achievements []UserAchievement `json:"achievements"`
}

// Detailed Achievement
// @Description Model of Achievement an user can have
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

// User Event
// @Description Model of Event that user have participated
type UserEvent struct {
	Name           string `json:"name"`
	WomID          string `json:"wom_id"`
	GuildID        string `json:"guild_id"`
	Placement      int16  `json:"placement"`
	PositionCutoff int16  `json:"position_cutoff"`
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
		}
	}

	return result
}

// Time teammates
// @Description Model that represents all teammates of a specific run
type TimeTeammates struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

// User Times
// @Description Model that represents all user times
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
	time := UserTime{
		Id:        0,
		Teammates: make([]TimeTeammates, 0),
	}

	for i := range rows {
		if rows[i].RunID != time.Id {
			if i != 0 {
				result = append(result, time)
			}
			time = UserTime{
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
		time.Teammates = append(time.Teammates, TimeTeammates{
			UserID:  rows[i].UserID,
			GuildID: rows[i].GuildID,
		})
	}

	// Add the final time object
	result = append(result, time)

	return result
}

// Users Model
// @Description Model of active guild members
type Users struct {
	Users []User `json:"users,omitempty"`
}

// TimeResponse Model
// @Description Return type of times endpoint retaining information to new and old time
type TimeResponse struct {
	BossName string `json:"boss_name"`
	Time     int    `json:"time"`
	OldTime  int    `json:"time_old"`
	RunID    int    `json:"run_id"`
}

// Time Model
// @Description Model of a fetched time and the team
type Time struct {
	Time     int        `json:"time"`
	BossName string     `json:"boss_name"`
	RunId    int        `json:"run_id"`
	Date     time.Time  `json:"date"`
	Team     []Teammate `json:"team"`
}

type Teammate struct {
	RunId   int    `json:"run_id"`
	GuildId string `json:"guild_id"`
	UserId  string `json:"user_id"`
}

type RSN struct {
	RSN     string `json:"rsn"`
	WomId   string `json:"wom_id"`
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
}

type DetailedEvent struct {
	Participations []EventParticipation `json:"participations"`
}

type EventParticipation struct {
	UserId    string `json:"user_id"`
	Placement int    `json:"placement"`
}

type Events struct {
	Events []database.Event `json:"events"`
}

type GuildTimes struct {
	guild_id         string                   `json:"guild_id"`
	pb_channel_id    string                   `json:"pb_channel_id"`
	bosses           []database.Boss          `json:"bosses"`
	categories       []database.Category      `json:"categories"`
	guild_bosses     []database.GuildBoss     `json:"guild_bosses"`
	guild_categories []database.GuildCategory `json:"guild_categories"`
	pbs              []Time                   `json:"pbs"`
	teammates        []User                   `json:"teammates"`
}

// Body Model
// @Description HTTP Body model for all responses
type Body struct {
	Content any `json:"content,omitempty"`
}

// Error Response Model
// @Description Model representing the error messages
type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Body Model
// @Description HTTP Body model for all responses
type Empty struct{}
