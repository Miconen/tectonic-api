// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Boss struct {
	Name        string
	DisplayName string
	Category    string
	Solo        bool
}

type Category struct {
	Thumbnail pgtype.Text
	Order     int16
	Name      string
}

type Guild struct {
	GuildID     string
	Multiplier  int32
	PbChannelID pgtype.Text
}

type GuildBoss struct {
	Boss    string
	GuildID string
	PbID    pgtype.Int4
}

type GuildCategory struct {
	GuildID   string
	Category  string
	MessageID pgtype.Text
}

type PointSource struct {
	GuildID string
	Source  string
	Points  int32
}

type Rsn struct {
	Rsn     string
	WomID   string
	UserID  string
	GuildID string
}

type Team struct {
	RunID   int32
	UserID  string
	GuildID string
}

type Time struct {
	Time     int32
	BossName string
	RunID    int32
	Date     pgtype.Timestamp
}

type User struct {
	UserID  string
	GuildID string
	Points  int32
}
