package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
)

type GetLeaderboardInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Limit   int32  `query:"limit" default:"50" minimum:"1" maximum:"1000" doc:"Maximum number of users to return"`
}
type GetLeaderboardOutput struct {
	Body []models.LeaderboardUser
}

func (s *Server) GetLeaderboard(ctx context.Context, input *GetLeaderboardInput) (*GetLeaderboardOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	} else if limit > 1000 {
		limit = 1000
	}

	params := database.GetLeaderboardParams{
		GuildID:   input.GuildID,
		UserLimit: limit,
	}

	logging.Get().DebugContext(ctx, "querying leaderboard from database", "guild_id", params.GuildID, "user_limit", params.UserLimit)
	rows, err := s.queries.GetLeaderboard(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if len(rows) == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}

	leaderboard := models.LeaderboardFromRows(rows)
	return &GetLeaderboardOutput{Body: leaderboard}, nil
}
