package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
)

type CompetitionResponse struct {
	Title            string                `json:"title"`
	ParticipantCount int                   `json:"participant_count"`
	Participants     []models.DetailedUser `json:"participants"`
	Accounts         []string              `json:"accounts"`
	Cutoff           int                   `json:"cutoff"`
	PointsGiven      int                   `json:"points_given"`
}

type EndCompetitionInput struct {
	GuildID       string `path:"guild_id" doc:"Guild Snowflake ID"`
	CompetitionID int    `path:"competition_id" doc:"WOM Competition ID"`
	Cutoff        int    `path:"cutoff" doc:"Minimum score cutoff"`
}
type EndCompetitionOutput struct {
	Body CompetitionResponse
}

func (s *Server) EndCompetition(ctx context.Context, input *EndCompetitionInput) (*EndCompetitionOutput, error) {
	c, err := s.womClient.GetCompetition(input.CompetitionID)
	if err != nil {
		return nil, s.womError(err)
	}

	emptyResponse := CompetitionResponse{
		Title:            c.Title,
		ParticipantCount: c.ParticipantCount,
		Participants:     []models.DetailedUser{},
		Accounts:         []string{},
		Cutoff:           input.Cutoff,
		PointsGiven:      0,
	}

	if len(c.Participations) == 0 {
		return &EndCompetitionOutput{Body: emptyResponse}, nil
	}

	rsns := make([]string, 0)
	accounts := make([]string, len(c.Participations))

	for i, val := range c.Participations {
		accounts[i] = val.Player.DisplayName
		if val.Progress.Gained < float64(input.Cutoff) {
			continue
		}
		rsns = append(rsns, val.Player.DisplayName)
	}

	if len(rsns) == 0 {
		emptyResponse.Accounts = accounts
		return &EndCompetitionOutput{Body: emptyResponse}, nil
	}

	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	userIDs, ei := database.WrapQuery(q.GetUserByRsn, ctx, rsns)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	points, err := q.UpdatePointsByEvent(ctx, database.UpdatePointsByEventParams{
		Event:   "event_participation",
		GuildID: input.GuildID,
		UserIds: userIDs,
	})
	if dbEi := database.ClassifyError(err); dbEi != nil {
		return nil, s.dbError(*dbEi)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}

	if len(points) == 0 {
		logging.Get().Info("no activated users found in competition")
	}

	users, ei := s.getDetailedUsers(ctx, userIDs, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	var pointsGiven int
	if len(points) > 0 {
		pointsGiven = int(points[0].GivenPoints)
	}

	return &EndCompetitionOutput{Body: CompetitionResponse{
		Title:            c.Title,
		ParticipantCount: c.ParticipantCount,
		Participants:     users,
		Accounts:         rsns,
		Cutoff:           input.Cutoff,
		PointsGiven:      pointsGiven,
	}}, nil
}

type CompetitionWinnersInput struct {
	GuildID       string `path:"guild_id" doc:"Guild Snowflake ID"`
	CompetitionID string `path:"competition_id" doc:"WOM Competition ID"`
}
type CompetitionWinnersOutput struct {
	Body any
}

func (s *Server) CompetitionWinners(ctx context.Context, input *CompetitionWinnersInput) (*CompetitionWinnersOutput, error) {
	// TODO: implement
	return nil, nil
}

type CompetitionTeamPositionInput struct {
	GuildID       string `path:"guild_id" doc:"Guild Snowflake ID"`
	CompetitionID string `path:"competition_id" doc:"WOM Competition ID"`
	Team          string `path:"team" doc:"Team name"`
}
type CompetitionTeamPositionOutput struct {
	Body any
}

func (s *Server) CompetitionTeamPosition(ctx context.Context, input *CompetitionTeamPositionInput) (*CompetitionTeamPositionOutput, error) {
	// TODO: implement
	return nil, nil
}
