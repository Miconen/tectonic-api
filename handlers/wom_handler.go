package handlers

import (
	"net/http"
	"strconv"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

type CompetitionResponse struct {
	Title            string                `json:"title"`
	ParticipantCount int                   `json:"participant_count"`
	Participants     []models.DetailedUser `json:"participants"`
	Accounts         []string              `json:"accounts"`
	Cutoff           int                   `json:"cutoff"`
	PointsGiven      int                   `json:"points_given"`
}

// @Summary		Handle Wise Old Man competitions
// @Description	Handle point giving and automatic data fetching through the Wise Old Man API
// @Tags			WOM
// @Accept			json
// @Produce		json
// @Param			guild_id		path		string	true	"Guild ID"
// @Param			competition_id	path		string	true	"Competition ID"
// @Param			cutoff			path		integer	true	"Cutoff"
// @Success		200				{object}	CompetitionResponse
// @Failure		400				{object}	models.Empty
// @Failure		401				{object}	models.Empty
// @Failure		409				{object}	models.Empty
// @Failure		429				{object}	models.Empty
// @Failure		500				{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/wom/competition/{competition_id}/cutoff/{cutoff} [GET]
func (s *Server) EndCompetition(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	id, err := strconv.Atoi(p["competition_id"])
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	cutoff, err := strconv.Atoi(p["cutoff"])
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	c, err := s.womClient.GetCompetition(id)
	if err != nil {
		s.handleWomError(err, jw)
		return
	}

	emptyResponse := CompetitionResponse{
		Title:            c.Title,
		ParticipantCount: c.ParticipantCount,
		Participants:     []models.DetailedUser{},
		Accounts:         []string{},
		Cutoff:           cutoff,
		PointsGiven:      0,
	}

	if len(c.Participations) == 0 {
		jw.WriteResponse(emptyResponse)
		return
	}

	rsns := make([]string, 0)
	accounts := make([]string, len(c.Participations))

	// Get participants valid for points
	for i, val := range c.Participations {
		accounts[i] = val.Player.DisplayName

		// Skip player if they didn't make the cutoff
		if val.Progress.Gained < float64(cutoff) {
			continue
		}

		rsns = append(rsns, val.Player.DisplayName)
	}

	if len(c.Participations) == 0 {
		jw.WriteResponse(emptyResponse)
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	q := s.queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	user_ids, ei := database.WrapQuery(q.GetUserByRsn, r.Context(), rsns)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	points_params := database.UpdatePointsByEventParams{
		Event:   "event_participation",
		GuildID: p["guild_id"],
		UserIds: user_ids,
	}

	points, err := q.UpdatePointsByEvent(r.Context(), points_params)
	if err != nil {
		ei = database.ClassifyError(err)
		s.handleDatabaseError(*ei, jw)
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	if len(points) == 0 {
		logging.Get().Info("no activated users found in competition")
	}

	users, ei := s.getDetailedUsers(r.Context(), user_ids, p["guild_id"])
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	var pointsGiven int
	if len(points) > 0 {
		pointsGiven = int(points[0].GivenPoints)
	}

	response := CompetitionResponse{
		Title:            c.Title,
		ParticipantCount: c.ParticipantCount,
		Participants:     users,
		Accounts:         rsns,
		Cutoff:           cutoff,
		PointsGiven:      pointsGiven,
	}

	jw.WriteResponse(response)
}

// @Summary		Event winners endpoint
// @Description	Handle past event winners and automatic data fetching through the Wise Old Man API
// @Tags			WOM
// @Accept			json
// @Produce		json
// @Param			guild_id		path		string	true	"Guild ID"
// @Param			competition_id	path		string	true	"Competition ID"
// @Param			cutoff			path		integer	true	"Cutoff"
// @Success		200				{object}	CompetitionResponse
// @Failure		400				{object}	models.ErrorResponse
// @Failure		401				{object}	models.ErrorResponse
// @Failure		409				{object}	models.ErrorResponse
// @Failure		429				{object}	models.ErrorResponse
// @Failure		500				{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/wom/winners/{competition_id} [GET]
func (s *Server) CompetitionWinners(w http.ResponseWriter, r *http.Request) {
}

// @Summary		Event winners by teamname endpoint
// @Description	Handle past event winners through teamname and automatic data fetching through the Wise Old Man API
// @Tags			WOM
// @Accept			json
// @Produce		json
// @Param			guild_id		path		string	true	"Guild ID"
// @Param			competition_id	path		string	true	"Competition ID"
// @Param			cutoff			path		integer	true	"Cutoff"
// @Success		200				{object}	CompetitionResponse
// @Failure		400				{object}	models.ErrorResponse
// @Failure		401				{object}	models.ErrorResponse
// @Failure		409				{object}	models.ErrorResponse
// @Failure		429				{object}	models.ErrorResponse
// @Failure		500				{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/wom/winners/{competition_id}/team/{team} [GET]
func (s *Server) CompetitionTeamPosition(w http.ResponseWriter, r *http.Request) {
}
