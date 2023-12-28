package handlers

// Guild Model
// @Description Model of guild data
type Guild struct {
	guild_id      string
	multiplier    int
	pb_channel_id string
}

// User Model
// @Description Model of active guild member
type User struct {
	user_id  string
	guild_id string
	points   int
}

// Users Model
// @Description Model of active guild members
type Users struct {
	users []User
}

// Time Model
// @Description Model of a fetched time and the team
type Time struct {
	time      int
	boss_name string
	run_id    int
	date      int
	team      []User
}

// Error Model
// @Description HTTP Error model
// @Description with content, error and code
type Error struct {
	content string
	error   string
	code    int
}
