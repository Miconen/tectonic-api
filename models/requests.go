package models

type InputGuild struct {
	GuildID DiscordSnowflake `json:"guild_id"`
}

type InputRecord struct {
	Value    int                `json:"value"      minimum:"1"`
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

type InputLegacyEvent struct {
	Name    string   `json:"name"     minLength:"1" maxLength:"128"`
	UserIDs []string `json:"user_ids" minItems:"1"`
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

type GuildCategoryMessage struct {
	MessageID DiscordSnowflake `json:"message_id"`
	Category  string           `json:"category"`
}

type PbUpdate struct {
	ChannelID        DiscordSnowflake       `json:"channel_id"`
	CategoryMessages []GuildCategoryMessage `json:"category_messages" minItems:"1"`
}

type UpdateGuildBody struct {
	Multiplier    *int              `json:"multiplier,omitempty" minimum:"1" maximum:"10"`
	ModChannelID  *DiscordSnowflake `json:"mod_channel_id,omitempty"`
	PbUpdate      *PbUpdate         `json:"pb_update,omitempty"`
	PositionCount *int              `json:"position_count,omitempty"`
}

type CreateGuildRankBody struct {
	Name         string  `json:"name"          minLength:"1" maxLength:"32"`
	MinPoints    int     `json:"min_points"    minimum:"0"`
	Icon         *string `json:"icon,omitempty"`
	RoleID       *string `json:"role_id,omitempty"`
	DisplayOrder int     `json:"display_order" minimum:"0"`
}

type UpdateGuildRankBody struct {
	MinPoints    *int    `json:"min_points,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	RoleID       *string `json:"role_id,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
}
