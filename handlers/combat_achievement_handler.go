package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/models"
)

type GetGuildCombatAchievementsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetGuildCombatAchievementsOutput struct {
	Body []database.GetGuildCombatAchievementsRow
}

func (s *Server) GetGuildCombatAchievements(ctx context.Context, input *GetGuildCombatAchievementsInput) (*GetGuildCombatAchievementsOutput, error) {
	cas, ei := database.WrapQuery(s.queries.GetGuildCombatAchievements, ctx, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetGuildCombatAchievementsOutput{Body: cas}, nil
}

type CreateCombatAchievementInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.CreateCombatAchievementBody
}

func (s *Server) CreateCombatAchievement(ctx context.Context, input *CreateCombatAchievementInput) (*struct{}, error) {
	err := s.queries.CreateCombatAchievement(ctx, database.CreateCombatAchievementParams{
		Name:        input.Body.Name,
		GuildID:     input.GuildID,
		PointSource: input.Body.PointSource,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type DeleteCombatAchievementInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	CAName  string `path:"ca_name" doc:"Combat Achievement Name"`
}

func (s *Server) DeleteCombatAchievement(ctx context.Context, input *DeleteCombatAchievementInput) (*struct{}, error) {
	rows, err := s.queries.DeleteCombatAchievement(ctx, database.DeleteCombatAchievementParams{
		Name:    input.CAName,
		GuildID: input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND)
	}
	return nil, nil
}

type CompleteCombatAchievementInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	CAName  string `path:"ca_name" doc:"Combat Achievement Name"`
	Body    models.CompleteCombatAchievementBody
}
type CompleteCombatAchievementOutput struct {
	Body []database.UpdatePointsByEventRow
}

func (s *Server) CompleteCombatAchievement(ctx context.Context, input *CompleteCombatAchievementInput) (*CompleteCombatAchievementOutput, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	ca, err := q.GetCombatAchievement(ctx, database.GetCombatAchievementParams{
		Name:    input.CAName,
		GuildID: input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	points, err := q.UpdatePointsByEvent(ctx, database.UpdatePointsByEventParams{
		Event:   ca.PointSource,
		GuildID: input.GuildID,
		UserIds: input.Body.UserIds,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	err = q.CompleteCombatAchievement(ctx, database.CompleteCombatAchievementParams{
		UserIds:               input.Body.UserIds,
		GuildID:               input.GuildID,
		CombatAchievementName: input.CAName,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return &CompleteCombatAchievementOutput{Body: points}, nil
}

type UserCombatAchievementInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
	CAName  string `path:"ca_name" doc:"Combat Achievement Name"`
}

func (s *Server) GiveUserCombatAchievement(ctx context.Context, input *UserCombatAchievementInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.GiveUserCombatAchievement, ctx, database.GiveUserCombatAchievementParams{
		UserID:                input.UserID,
		GuildID:               input.GuildID,
		CombatAchievementName: input.CAName,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

func (s *Server) RemoveUserCombatAchievement(ctx context.Context, input *UserCombatAchievementInput) (*struct{}, error) {
	rows, err := s.queries.RemoveUserCombatAchievement(ctx, database.RemoveUserCombatAchievementParams{
		UserID:                input.UserID,
		GuildID:               input.GuildID,
		CombatAchievementName: input.CAName,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_COMBAT_ACHIEVEMENT_NOT_FOUND)
	}
	return nil, nil
}
