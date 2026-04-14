package models

type InputGuild struct {
	GuildID DiscordSnowflake `json:"guild_id"`
}

type InputTime struct {
	Time     int                `json:"time"      minimum:"1"`
	BossName string             `json:"boss_name"  minLength:"1" maxLength:"50"`
	UserIDs  []DiscordSnowflake `json:"user_ids"   minItems:"1"  maxItems:"8"`
}

type InputTeammate struct {
	UserID  DiscordSnowflake `json:"user_id"`
	GuildID DiscordSnowflake `json:"guild_id"`
}

type InputEvent struct {
	EventID        int      `json:"event_id"         minimum:"1"`
	TeamNames      []string `json:"team_names,omitempty"`
	PositionCutoff int      `json:"position_cutoff,omitempty" minimum:"1" maximum:"3"`
}

type CreateUserBody struct {
	UserID DiscordSnowflake `json:"user_id"`
	RSN    RSN              `json:"rsn"`
}

type CreateRsnBody struct {
	RSN RSN `json:"rsn"`
}

type CreateCombatAchievementBody struct {
	Name        string `json:"name"         minLength:"1" maxLength:"64"`
	PointSource string `json:"point_source" minLength:"1" maxLength:"32"`
}

type CompleteCombatAchievementBody struct {
	UserIDs []DiscordSnowflake `json:"user_ids" minItems:"1" maxItems:"8"`
}
