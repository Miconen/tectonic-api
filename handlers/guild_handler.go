package handlers

import (
	"context"
	"math/big"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetGuildInput struct {
	GuildID models.DiscordSnowflake `path:"guild_id" doc:"Guild Snowflake ID"`
}

type GetGuildOutput struct {
	Body database.Guild
}

func (s *Server) GetGuild(ctx context.Context, input *GetGuildInput) (*GetGuildOutput, error) {
	guild, err := s.queries.GetGuild(ctx, string(input.GuildID))
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
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
	Body    struct {
		Multiplier   *float64                 `json:"multiplier,omitempty"`
		ModChannelID *models.DiscordSnowflake `json:"mod_channel_id,omitempty"`
		PbUpdate     *models.PbUpdate         `json:"pb_update,omitempty"`
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

	// Handle PB update (channel + category messages)
	var pbChannelID string
	if input.Body.PbUpdate != nil {
		pbChannelID = input.Body.PbUpdate.ChannelID.String()

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

	// Convert *float64 to pgtype.Numeric
	var multiplier pgtype.Numeric
	if input.Body.Multiplier != nil {
		multiplier.Valid = true
		multiplier.Int = big.NewInt(int64(*input.Body.Multiplier * 100))
		multiplier.Exp = -2
	}

	// Convert *DiscordSnowflake to string
	var modChannelID string
	if input.Body.ModChannelID != nil {
		modChannelID = utils.DerefOr(input.Body.ModChannelID, "").String()
	}

	_, err = q.UpdateGuild(ctx, database.UpdateGuildParams{
		Multiplier:   multiplier,
		PbChannelID:  pbChannelID,
		ModChannelID: modChannelID,
		GuildID:      input.GuildID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	tx.Commit(ctx)
	return nil, nil
}
