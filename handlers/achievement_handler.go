package handlers

import (
	"context"

	"tectonic-api/database"
)

type (
	GetAchievementsInput  struct{}
	GetAchievementsOutput struct {
		Body []database.Achievement
	}
)

func (s *Server) GetAchievements(ctx context.Context, input *GetAchievementsInput) (*GetAchievementsOutput, error) {
	achievements, err := s.queries.GetAchievements(ctx)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetAchievementsOutput{Body: achievements}, nil
}

// Shared input types

type AchievementByIdInput struct {
	GuildID     string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID      string `path:"user_id" doc:"User Snowflake ID"`
	Achievement string `path:"achievement" doc:"Achievement name"`
}

type AchievementByRsnInput struct {
	GuildID     string `path:"guild_id" doc:"Guild Snowflake ID"`
	RSN         string `path:"rsn" doc:"RuneScape Name"`
	Achievement string `path:"achievement" doc:"Achievement name"`
}

// Give

func (s *Server) GiveAchievementById(ctx context.Context, input *AchievementByIdInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.GiveAchievementById, ctx, database.GiveAchievementByIdParams{
		UserID:          input.UserID,
		AchievementName: input.Achievement,
		GuildID:         input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

func (s *Server) GiveAchievementByRsn(ctx context.Context, input *AchievementByRsnInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.GiveAchievementByRsn, ctx, database.GiveAchievementByRsnParams{
		Rsn:             input.RSN,
		AchievementName: input.Achievement,
		GuildID:         input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

// Remove

func (s *Server) RemoveAchievementById(ctx context.Context, input *AchievementByIdInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.RemoveAchievementById, ctx, database.RemoveAchievementByIdParams{
		UserID:          input.UserID,
		AchievementName: input.Achievement,
		GuildID:         input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

func (s *Server) RemoveAchievementByRsn(ctx context.Context, input *AchievementByRsnInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.RemoveAchievementByRsn, ctx, database.RemoveAchievementByRsnParams{
		Rsn:             input.RSN,
		AchievementName: input.Achievement,
		GuildID:         input.GuildID,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}
