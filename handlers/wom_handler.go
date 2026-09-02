package handlers

import (
	"context"
	"strconv"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
)

type CompetitionResponse struct {
	Title            string                `json:"title"`
	ParticipantCount int                   `json:"participant_count"`
	Participants     []models.DetailedUser `json:"participants"`
	Accounts         []models.UserRsn      `json:"accounts"`
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

func (s *Server) EndCompetition(
	ctx context.Context,
	input *EndCompetitionInput,
) (*EndCompetitionOutput, error) {
	competition, err := s.womClient.GetCompetition(input.CompetitionID)
	if err != nil {
		return nil, s.womError(err, models.ERROR_WOM_COMPETITION_NOT_FOUND)
	}

	emptyResponse := CompetitionResponse{
		Title:            competition.Title,
		ParticipantCount: competition.ParticipantCount,
		Participants:     []models.DetailedUser{},
		Accounts:         []models.UserRsn{},
		Cutoff:           input.Cutoff,
		PointsGiven:      0,
	}

	accounts := make([]models.UserRsn, 0, len(competition.Participations))

	for _, participation := range competition.Participations {
		if participation.Progress.Gained < float64(input.Cutoff) {
			continue
		}

		accounts = append(accounts, models.UserRsn{
			RSN:   participation.Player.DisplayName,
			WomId: strconv.Itoa(participation.PlayerID),
		})
	}

	if len(accounts) == 0 {
		return &EndCompetitionOutput{Body: emptyResponse}, nil
	}

	womIDs := make([]string, len(accounts))
	for i, account := range accounts {
		womIDs[i] = account.WomId
	}

	tx, err := database.CreateTx(ctx)
	if err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	userIDs, ei := database.WrapQuery(
		q.GetGuildUserByWom,
		ctx,
		database.GetGuildUserByWomParams{
			GuildID: input.GuildID,
			WomIds:  womIDs,
		},
	)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	points, err := q.UpdatePointsByEvent(
		ctx,
		database.UpdatePointsByEventParams{
			Event:   "event_participation",
			GuildID: input.GuildID,
			UserIds: userIDs,
		},
	)
	if dbEi := database.ClassifyError(err); dbEi != nil {
		return nil, s.dbError(*dbEi)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}

	if len(points) == 0 {
		logging.Get().Info(
			"no activated users found in competition",
			"guild_id", input.GuildID,
			"competition_id", input.CompetitionID,
		)

		return &EndCompetitionOutput{Body: CompetitionResponse{
			Title:            competition.Title,
			ParticipantCount: competition.ParticipantCount,
			Participants:     []models.DetailedUser{},
			Accounts:         accounts,
			Cutoff:           input.Cutoff,
			PointsGiven:      0,
		}}, nil
	}

	updatedUserIDs := make([]string, len(points))
	for i, point := range points {
		updatedUserIDs[i] = point.UserID
	}

	users, ei := s.getDetailedUsers(ctx, updatedUserIDs, input.GuildID)
	if ei != nil {
		return nil, s.dbError(*ei)
	}

	return &EndCompetitionOutput{Body: CompetitionResponse{
		Title:            competition.Title,
		ParticipantCount: competition.ParticipantCount,
		Participants:     users,
		Accounts:         accounts,
		Cutoff:           input.Cutoff,
		PointsGiven:      int(points[0].GivenPoints),
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
