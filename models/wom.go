package models

import "time"

type WomCompetition struct {
	ID               int              `json:"id"`
	Title            string           `json:"title"`
	Metric           string           `json:"metric"`
	Type             string           `json:"type"`
	StartsAt         time.Time        `json:"startsAt"`
	EndsAt           time.Time        `json:"endsAt"`
	GroupID          int              `json:"groupId"`
	Score            int              `json:"score"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	ParticipantCount int              `json:"participantCount"`
	Participations   []Participations `json:"participations"`
}
type Player struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"displayName"`
	Type           string    `json:"type"`
	Build          string    `json:"build"`
	Country        any       `json:"country"`
	Status         string    `json:"status"`
	Patron         bool      `json:"patron"`
	Exp            int       `json:"exp"`
	Ehp            float64   `json:"ehp"`
	Ehb            float64   `json:"ehb"`
	Ttm            float64   `json:"ttm"`
	Tt200M         float64   `json:"tt200m"`
	RegisteredAt   time.Time `json:"registeredAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	LastChangedAt  time.Time `json:"lastChangedAt"`
	LastImportedAt any       `json:"lastImportedAt"`
}
type Progress struct {
	Gained float64 `json:"gained"`
	Start  float64 `json:"start"`
	End    float64 `json:"end"`
}
type Levels struct {
	Gained int `json:"gained"`
	Start  int `json:"start"`
	End    int `json:"end"`
}
type Participations struct {
	PlayerID      int       `json:"playerId"`
	CompetitionID int       `json:"competitionId"`
	TeamName      string    `json:"teamName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Player        Player    `json:"player"`
	Progress      Progress  `json:"progress"`
	Levels        Levels    `json:"levels"`
}
