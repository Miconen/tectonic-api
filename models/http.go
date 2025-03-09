package models

import "time"

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

// Body Model
// @Description HTTP Body model for all responses
type Body struct {
	Content interface{} `json:"content,omitempty"`
}

// Body Model
// @Description HTTP Body model for all responses
type Empty struct{}
