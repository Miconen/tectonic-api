package models

import (
	"tectonic-api/database"
	"time"
)

// InputGuild Model
// @Description Model of new guild data
type InputGuild struct {
	GuildId     string `json:"guild_id"`
	Multiplier  int    `json:"multiplier" optional:"true"`
	PbChannelId string `json:"pb_channel_id" optional:"true"`
}

// UpdateGuild Model
// @Description Model of updated guild data (Exists so PbChannelId can be properly parsed)
type UpdateGuild struct {
	GuildId     string `json:"guild_id"`
	Multiplier  int    `json:"multiplier"`
	PbChannelId string `json:"pb_channel_id"`
}

// Guild Model
// @Description Model of guild data
type Guild struct {
	GuildId     string     `json:"guild_id"`
	Multiplier  int        `json:"multiplier"`
	PbChannelId NullString `json:"pb_channel_id"`
}

// InputUser Model
// @Description Model of new active guild member
type InputUser struct {
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
	RSN     string `json:"rsn"`
}

type GuildUsers interface {
	GetUserIDs() []string
	GetGuildID() string
}

// User Model
// @Description Model of active guild member
type User struct {
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
	Points  int    `json:"points"`
}

// Users Model
// @Description Model of active guild members
type Users struct {
	Users []User `json:"users,omitempty"`
}

// InputTime Model
// @Description Model of a new time
type InputTime struct {
	GuildId  string   `json:"guild_id"`
	Time     int      `json:"time"`
	BossName string   `json:"boss_name"`
	UserIds  []string `json:"user_ids"`
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

type InputRSN struct {
	RSN     string `json:"rsn"`
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
}

type RSN struct {
	RSN     string `json:"rsn"`
	WomId   string `json:"wom_id"`
	UserId  string `json:"user_id"`
	GuildId string `json:"guild_id"`
}

type GuildTimes struct {
	guild_id         string          `json:"guild_id"`
	pb_channel_id    string          `json:"pb_channel_id"`
	bosses           []database.Boss          `json:"bosses"`
	categories       []database.Category      `json:"categories"`
	guild_bosses     []database.GuildBoss     `json:"guild_bosses"`
	guild_categories []database.GuildCategory `json:"guild_categories"`
	pbs              []Time          `json:"pbs"`
	teammates        []User          `json:"teammates"`
}

// Body Model
// @Description HTTP Body model for all responses
type Body struct {
	Content interface{} `json:"content,omitempty"`
}

// Body Model
// @Description HTTP Body model for all responses
type Empty struct{}
