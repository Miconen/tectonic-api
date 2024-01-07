package models

// Guild Model
// @Description Model of guild data
type Guild struct {
	GuildId     string `json:"guild_id"`
	Multiplier  int    `json:"multiplier"`
	PbChannelId string `json:"pb_channel_id"`
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
	Users []User `json:"users"`
}

// Time Model
// @Description Model of a fetched time and the team
type Time struct {
	Time     int    `json:"time"`
	BossName string `json:"boss_name"`
	RunId    int    `json:"run_id"`
	Date     int    `json:"date"`
	Team     Users  `json:"team"`
}

// Error Model
// @Description HTTP Error model with content, error and code
type Error struct {
	Content string `json:"content"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}
