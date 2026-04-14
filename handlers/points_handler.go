package handlers

import (
	"context"
	"strings"

	"tectonic-api/database"
	"tectonic-api/models"
)

type UpdatePointsInput struct {
	GuildID    string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserIDs    string `path:"user_ids" doc:"Comma-separated User Snowflake IDs"`
	PointEvent string `path:"point_event" doc:"Point event name"`
}
type UpdatePointsOutput struct {
	Body any
}

func (s *Server) UpdatePoints(ctx context.Context, input *UpdatePointsInput) (*UpdatePointsOutput, error) {
	params := database.UpdatePointsByEventParams{
		Event:   input.PointEvent,
		GuildID: input.GuildID,
		UserIds: strings.Split(input.UserIDs, ","),
	}

	user, err := s.queries.UpdatePointsByEvent(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		if ei.Recoverable && ei.Code == "23502" {
			return nil, models.NewTectonicError(models.ERROR_POINT_SOURCE_NOT_FOUND)
		}
		return nil, s.dbError(*ei)
	}

	if len(user) == 0 {
		return nil, models.NewTectonicError(models.ERROR_POINT_SOURCE_NOT_FOUND)
	}
	return &UpdatePointsOutput{Body: user}, nil
}

type UpdatePointsCustomInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	UserIDs string `path:"user_ids" doc:"Comma-separated User Snowflake IDs"`
	Points  int    `path:"points" doc:"Points to add"`
}
type UpdatePointsCustomOutput struct {
	Body any
}

func (s *Server) UpdatePointsCustom(ctx context.Context, input *UpdatePointsCustomInput) (*UpdatePointsCustomOutput, error) {
	params := database.UpdatePointsCustomParams{
		Points:  int32(input.Points),
		UserIds: strings.Split(input.UserIDs, ","),
		GuildID: input.GuildID,
	}

	user, err := s.queries.UpdatePointsCustom(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if len(user) == 0 {
		return nil, models.NewTectonicError(models.ERROR_POINT_SOURCE_NOT_FOUND)
	}
	return &UpdatePointsCustomOutput{Body: user}, nil
}

type GetPointSourcesInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetPointSourcesOutput struct {
	Body any
}

func (s *Server) GetPointSources(ctx context.Context, input *GetPointSourcesInput) (*GetPointSourcesOutput, error) {
	events, ei := database.WrapQuery(s.queries.GetGuildPointSources, ctx, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return &GetPointSourcesOutput{Body: events}, nil
}

type UpdateGuildPointSourceInput struct {
	GuildID     string `path:"guild_id" doc:"Guild Snowflake ID"`
	PointSource string `path:"point_source" doc:"Point source name"`
	Points      int    `path:"points" doc:"New point value"`
}

func (s *Server) UpdateGuildPointSource(ctx context.Context, input *UpdateGuildPointSourceInput) (*struct{}, error) {
	params := database.UpdateGuildPointSourceParams{
		Points:      int32(input.Points),
		GuildID:     input.GuildID,
		PointSource: input.PointSource,
	}

	rowsAf, err := s.queries.UpdateGuildPointSource(ctx, params)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	if rowsAf == 0 {
		return nil, models.NewTectonicError(models.ERROR_POINT_SOURCE_NOT_FOUND)
	}
	return nil, nil
}
