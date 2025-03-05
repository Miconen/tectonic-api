package database

type UserData struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	RSNs    []Rsn  `json:"rsns"`
}

type DetailedTime struct {
	Time      Time       `json:"time"`
	Teammates []UserData `json:"team"`
}

type DetailedUser struct {
	User  UserData       `json:"user"`
	Times []DetailedTime `json:"times"`
}
