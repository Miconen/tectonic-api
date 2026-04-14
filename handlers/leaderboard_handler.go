package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
)

type GetLeaderboardInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetLeaderboardOutput struct {
	Body []database.UserData
}

func (s *Server) GetLeaderboard(ctx context.Context, input *GetLeaderboardInput) (*GetLeaderboardOutput, error) {
	params := database.GetLeaderboardParams{
		GuildID:   input.GuildID,
		UserLimit: 50,
	}

	logging.Get().DebugContext(ctx, "querying leaderboard from database", "guild_id", params.GuildID, "user_limit", params.UserLimit)
	rows, err := s.queries.GetLeaderboard(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if len(rows) == 0 {
		return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
	}

	leaderboard := database.NewLeaderboardFromRows(rows)
	return &GetLeaderboardOutput{Body: leaderboard}, nil
}
