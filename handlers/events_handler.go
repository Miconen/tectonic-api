package handlers

import (
	"context"
	"fmt"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetEventsInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
}
type GetEventsOutput struct {
	Body any
}

func (s *Server) GetEvents(ctx context.Context, input *GetEventsInput) (*GetEventsOutput, error) {
	events, ei := database.WrapQuery(s.queries.GetGuildEvents, ctx, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	if len(events) == 0 {
		return &GetEventsOutput{Body: []int{}}, nil
	}
	return &GetEventsOutput{Body: events}, nil
}

type GetDetailedEventInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	EventID string `path:"event_id" doc:"Event WOM ID"`
}
type GetDetailedEventOutput struct {
	Body models.DetailedEvent
}

func (s *Server) GetDetailedEvent(ctx context.Context, input *GetDetailedEventInput) (*GetDetailedEventOutput, error) {
	events, ei := database.WrapQuery(s.queries.GetEventParticipation, ctx, input.EventID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	res := models.DetailedEvent{
		Participations: utils.MapField(events, func(p database.GetEventParticipationRow) models.EventParticipation {
			return models.EventParticipation{
				UserId:    p.UserID,
				Placement: int(p.Placement),
			}
		}),
	}
	return &GetDetailedEventOutput{Body: res}, nil
}

type RegisterEventInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	Body    models.InputEvent
}

func (s *Server) RegisterEvent(ctx context.Context, input *RegisterEventInput) (*struct{}, error) {
	if input.Body.PositionCutoff == 0 {
		input.Body.PositionCutoff = 3
	}

	c, err := s.womClient.GetCompetition(input.Body.EventID)
	if err != nil {
		return nil, s.womError(err)
	}

	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	solo := c.Type != "team"

	if len(input.Body.TeamNames) > 0 {
		input.Body.PositionCutoff = len(input.Body.TeamNames)
	}

	q := s.queries.WithTx(tx)
	ei := database.WrapExec(q.CreateEvent, ctx, database.CreateEventParams{
		Name:           c.Title,
		WomID:          fmt.Sprintf("%d", c.ID),
		GuildID:        input.GuildID,
		PositionCutoff: int16(input.Body.PositionCutoff),
		Solo:           solo,
	})
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	if c.Type == "classic" {
		winners := utils.MapField(c.Participations, func(p models.Participations) string {
			return fmt.Sprintf("%d", p.PlayerID)
		})[0:3]

		ei = database.WrapExec(q.InsertEventParticipants, ctx, database.InsertEventParticipantsParams{
			ParticipantIds: winners,
			GuildID:        input.GuildID,
			WomID:          fmt.Sprintf("%d", c.ID),
		})
		if ei != nil {
			return nil, s.dbError(*ei)
		}
	} else if c.Type == "team" && len(input.Body.TeamNames) != 0 {
		posMap := make(map[string]int32)
		for i := range input.Body.TeamNames {
			posMap[input.Body.TeamNames[i]] = int32(i + 1)
		}

		var ids []string
		var positions []int32
		for _, participation := range c.Participations {
			if position, ok := posMap[participation.TeamName]; ok {
				ids = append(ids, fmt.Sprintf("%d", participation.PlayerID))
				positions = append(positions, position)
			}
		}

		if len(ids) > 0 {
			ei = database.WrapExec(q.InsertEventTeams, ctx, database.InsertEventTeamsParams{
				ParticipantIds:        ids,
				ParticipantPlacements: positions,
				GuildID:               input.GuildID,
				WomID:                 fmt.Sprintf("%d", c.ID),
			})
			if ei != nil {
				return nil, s.dbError(*ei)
			}
		}
	} else {
		return nil, models.NewTectonicError(models.ERROR_WRONG_BODY)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	return nil, nil
}

type DeleteEventInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	EventID string `path:"event_id" doc:"Event WOM ID"`
}

func (s *Server) DeleteEvent(ctx context.Context, input *DeleteEventInput) (*struct{}, error) {
	ei := database.WrapExec(s.queries.DeleteEvent, ctx, input.EventID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}
	return nil, nil
}

type UpdateEventInput struct {
	GuildID string `path:"guild_id" doc:"Guild Snowflake ID"`
	EventID string `path:"event_id" doc:"Event WOM ID"`
	Body    struct {
		Name           string         `json:"name"`
		PositionCutoff pgtype.Numeric `json:"position_cutoff"`
		Solo           bool           `json:"solo"`
	}
}

func (s *Server) UpdateEvent(ctx context.Context, input *UpdateEventInput) (*struct{}, error) {
	tx, err := database.CreateTx(ctx)
	if err != nil {
		logging.Get().Error("Error creating transaction", "error", err)
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	_, err = q.UpdateEvent(ctx, database.UpdateEventParams{
		Name:           input.Body.Name,
		PositionCutoff: input.Body.PositionCutoff,
		Solo:           input.Body.Solo,
		GuildID:        input.GuildID,
		WomID:          input.EventID,
	})
	if ei := database.ClassifyError(err); ei != nil {
		return nil, s.dbError(*ei)
	}

	tx.Commit(ctx)
	return nil, nil
}
