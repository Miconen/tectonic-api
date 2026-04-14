package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetGuildInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetGuildOutput struct {
	Body any
}

func (s *Server) GetGuild(ctx context.Context, input *GetGuildInput) (*GetGuildOutput, error) {
	guild, err := s.queries.GetGuild(ctx, input.GuildID)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetGuildOutput{Body: guild}, nil
}

type CreateGuildInput struct {
	Body models.InputGuild
}

func (s *Server) CreateGuild(ctx context.Context, input *CreateGuildInput) (*struct{}, error) {
	_, err := s.queries.CreateGuild(ctx, input.Body.GuildId)
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

type GuildCategoryMessage struct {
	MessageID string `json:"message_id"`
	Category  string `json:"category"`
}

type UpdateGuildInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    struct {
		Multiplier       pgtype.Numeric         `json:"multiplier"`
		PbChannelID      string                 `json:"pb_channel_id"`
		ModChannelID     string                 `json:"mod_channel_id"`
		CategoryMessages []GuildCategoryMessage `json:"category_messages"`
	}
}

func (s *Server) UpdateGuild(ctx context.Context, input *UpdateGuildInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	categories := make([]string, len(input.Body.CategoryMessages))
	messageIds := make([]string, len(input.Body.CategoryMessages))
	for i, v := range input.Body.CategoryMessages {
		categories[i] = v.Category
		messageIds[i] = v.MessageID
	}

	if len(input.Body.CategoryMessages) > 0 {
		params := database.UpdateCategoryMessageIdsParams{
			GuildID:    input.GuildID,
			Categories: categories,
			MessageIds: messageIds,
		}
		_, err := q.UpdateCategoryMessageIds(ctx, params)
		if ei := database.ClassifyError(err); ei != nil {
			return nil, s.dbError(*ei)
		}
	}

	guildParams := database.UpdateGuildParams{
		Multiplier:   input.Body.Multiplier,
		PbChannelID:  input.Body.PbChannelID,
		ModChannelID: input.Body.ModChannelID,
		GuildID:      input.GuildID,
	}

	_, err = q.UpdateGuild(ctx, guildParams)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	tx.Commit(ctx)
	return nil, nil
}
