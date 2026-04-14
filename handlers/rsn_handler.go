package handlers

import (
	"context"
	"strconv"

	"tectonic-api/database"
	"tectonic-api/models"
)

type CreateRsnBody struct {
	RSN string `json:"rsn"`
}

type CreateRSNInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
	Body    CreateRsnBody
}

func (s *Server) CreateRSN(ctx context.Context, input *CreateRSNInput) (*struct{}, error) {
	wom, err := s.womClient.GetWom(input.Body.RSN)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_RSN_NOT_FOUND)
	}

	params := database.CreateRsnParams{
		GuildID: input.GuildID,
		UserID:  input.UserID,
		WomID:   strconv.Itoa(wom.Id),
		Rsn:     wom.DisplayName,
	}

	err = s.queries.CreateRsn(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		if ei.Recoverable {
			switch ei.Code {
			case "23503":
				return nil, models.NewTectonicError(models.ERROR_USER_NOT_FOUND)
			case "23505":
				return nil, models.NewTectonicError(models.ERROR_RSN_EXISTS)
			}
		}
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type RemoveRSNInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserID  string `path:"user_id" doc:"User Snowflake ID"`
	RSN     string `path:"rsn" doc:"RuneScape Name"`
}

func (s *Server) RemoveRSN(ctx context.Context, input *RemoveRSNInput) (*struct{}, error) {
	rows, err := s.queries.DeleteRsn(ctx, database.DeleteRsnParams{
		GuildID: input.GuildID,
		UserID:  input.UserID,
		Rsn:     input.RSN,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}
	if rows == 0 {
		return nil, models.NewTectonicError(models.ERROR_RSN_NOT_FOUND)
	}
	return nil, nil
}
