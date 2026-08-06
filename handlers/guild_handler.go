package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
)

type GetGuildInput struct {
	GuildID  models.DiscordSnowflake `path:"guild_id" doc:"Guild Snowflake ID"`
	Detailed bool                    `query:"detailed" default:"false" doc:"Fetch detailed guild information"`
}

type GetGuildOutput struct {
	Body models.GuildResponse
}

func (s *Server) GetGuild(ctx context.Context, input *GetGuildInput) (*GetGuildOutput, error) {
	if input.Detailed {
		row, err := s.queries.GetDetailedGuild(ctx, string(input.GuildID))
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
		guild := models.GuildResponseFromDetailedRow(row)
		return &GetGuildOutput{Body: guild}, nil
	}

	row, err := s.queries.GetGuild(ctx, string(input.GuildID))
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	guild := models.GuildResponseFromRow(row)
	return &GetGuildOutput{Body: guild}, nil
}

type CreateGuildInput struct {
	Body models.InputGuild
}

func (s *Server) CreateGuild(ctx context.Context, input *CreateGuildInput) (*struct{}, error) {
	_, err := s.queries.CreateGuild(ctx, string(input.Body.GuildID))
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type DeleteGuildInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}

func (s *Server) DeleteGuild(ctx context.Context, input *DeleteGuildInput) (*struct{}, error) {
	rows, err := s.queries.DeleteGuild(ctx, input.GuildID)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_GUILD_NOT_FOUND)
	}
	return nil, nil
}

type UpdateGuildInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.UpdateGuildBody
}

func (s *Server) UpdateGuild(ctx context.Context, input *UpdateGuildInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	// Handle PB update (channel + category messages)
	var pbChannelID *string
	if input.Body.PbUpdate != nil {
		pbChannelID = utils.Ptr(input.Body.PbUpdate.ChannelID.String())

		categories := make([]string, len(input.Body.PbUpdate.CategoryMessages))
		messageIds := make([]string, len(input.Body.PbUpdate.CategoryMessages))
		for i, v := range input.Body.PbUpdate.CategoryMessages {
			categories[i] = v.Category
			messageIds[i] = v.MessageID.String()
		}

		_, err := q.UpdateCategoryMessageIds(ctx, database.UpdateCategoryMessageIdsParams{
			GuildID:    input.GuildID,
			Categories: categories,
			MessageIds: messageIds,
		})
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
	}

	_, err = q.UpdateGuild(ctx, database.UpdateGuildParams{
		Multiplier:    input.Body.Multiplier,
		PbChannelID:   pbChannelID,
		ModChannelID:  input.Body.ModChannelID.PtrString(),
		LogChannelID:  input.Body.LogChannelID.PtrString(),
		PositionCount: input.Body.PositionCount,
		GuildID:       input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	tx.Commit(ctx)
	return nil, nil
}
