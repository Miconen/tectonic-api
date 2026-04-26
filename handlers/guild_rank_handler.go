package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetGuildRanksInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetGuildRanksOutput struct {
	Body []models.GuildRankResponse
}

func (s *Server) GetGuildRanks(ctx context.Context, input *GetGuildRanksInput) (*GetGuildRanksOutput, error) {
	rows, ei := database.WrapQuery(s.queries.GetGuildRanks, ctx, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetGuildRanksOutput{Body: models.GuildRanksFromRows(rows)}, nil
}

type CreateGuildRankInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.CreateGuildRankBody
}

func (s *Server) CreateGuildRank(ctx context.Context, input *CreateGuildRankInput) (*struct{}, error) {
	var icon pgtype.Text
	if input.Body.Icon != nil {
		icon = pgtype.Text{String: *input.Body.Icon, Valid: true}
	}

	var roleID pgtype.Text
	if input.Body.RoleID != nil {
		roleID = pgtype.Text{String: *input.Body.RoleID, Valid: true}
	}

	params := database.CreateGuildRankParams{
		GuildID:      input.GuildID,
		Name:         input.Body.Name,
		MinPoints:    int32(input.Body.MinPoints),
		Icon:         icon,
		RoleID:       roleID,
		DisplayOrder: int16(input.Body.DisplayOrder),
	}
	ei := database.WrapExec(s.queries.CreateGuildRank, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type UpdateGuildRankInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Name    string `path:"name" doc:"Rank name"`
	Body    models.UpdateGuildRankBody
}

func (s *Server) UpdateGuildRank(ctx context.Context, input *UpdateGuildRankInput) (*struct{}, error) {
	var minPoints int32
	if input.Body.MinPoints != nil {
		minPoints = int32(*input.Body.MinPoints)
	}

	var icon string
	if input.Body.Icon != nil {
		icon = *input.Body.Icon
	}

	var roleID string
	if input.Body.RoleID != nil {
		roleID = *input.Body.RoleID
	}

	var displayOrder int16
	if input.Body.DisplayOrder != nil {
		displayOrder = int16(*input.Body.DisplayOrder)
	}

	params := database.UpdateGuildRankParams{
		GuildID:      input.GuildID,
		Name:         input.Name,
		MinPoints:    minPoints,
		Icon:         icon,
		RoleID:       roleID,
		DisplayOrder: displayOrder,
	}
	rows, ei := database.WrapQuery(s.queries.UpdateGuildRank, ctx, params)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_GUILD_RANK_NOT_FOUND)
	}
	return nil, nil
}

type DeleteGuildRankInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Name    string `path:"name" doc:"Rank name"`
}

func (s *Server) DeleteGuildRank(ctx context.Context, input *DeleteGuildRankInput) (*struct{}, error) {
	rows, ei := database.WrapQuery(s.queries.DeleteGuildRank, ctx, database.DeleteGuildRankParams{
		GuildID: input.GuildID,
		Name:    input.Name,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_GUILD_RANK_NOT_FOUND)
	}
	return nil, nil
}
