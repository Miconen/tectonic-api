package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

type CompetitionResponse struct {
	Title            string                      `json:"title"`
	ParticipantCount int                         `json:"participant_count"`
	Participants     []database.DetailedUserJSON `json:"participants"`
	Accounts         []string                    `json:"accounts"`
	Cutoff           int                         `json:"cutoff"`
	PointsGiven      int                         `json:"points_given"`
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
func EndCompetition(w http.ResponseWriter, r *http.Request) {
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

	competition, err := utils.GetCompetition(id)
	if err != nil {
		// TODO: differentiate response errors from request errors
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	if len(competition.Participations) == 0 {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	params := database.GetDetailedUsersByRSNParams{
		GuildID: p["guild_id"],
		Rsns:    []string{},
	}

	accounts := make([]string, len(competition.Participations))

	// Get participants valid for points
	for i, val := range competition.Participations {
		accounts[i] = val.Player.DisplayName

		// Skip player if they didn't make the cutoff
		if val.Progress.Gained < float64(cutoff) {
			continue
		}

		params.Rsns = append(params.Rsns, val.Player.DisplayName)
	}

	if len(params.Rsns) == 0 {
		jw.WriteError(models.ERROR_PARTICIPATION_NOT_FOUND)
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	rows, err := q.GetDetailedUsersByRSN(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	users := make([]database.DetailedUserJSON, 0, len(rows))
	for _, row := range rows {
		user := database.DetailedUserJSON{UserID: row.UserID, GuildID: row.GuildID, Points: row.Points, RSNs: row.Rsns, Times: row.Times}
		users = append(users, user)
	}

	user_ids := make([]string, len(users))
	for i, v := range users {
		user_ids[i] = v.UserID
	}

	given_params := database.GetPointsValueParams{
		Event:   "event_participation",
		GuildID: p["guild_id"],
	}

	given, err := q.GetPointsValue(r.Context(), given_params)
	ei = database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_API_UNAVAILABLE)
		return
	}

	points_params := database.UpdatePointsByEventParams{
		Event:   "event_participation",
		GuildID: p["guild_id"],
		UserIds: user_ids,
	}

	points, err := q.UpdatePointsByEvent(r.Context(), points_params)
	ei = database.ClassifyError(err)
	if err != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	tx.Commit(r.Context())

	if len(points) == 0 {
		log.Info("no activated users found in competition")
	}

	for i := range users {
		// Users stored points dont include given points, so we add them here
		users[i].Points = given
	}

	response := CompetitionResponse{
		Title:            competition.Title,
		ParticipantCount: competition.ParticipantCount,
		Participants:     users,
		Accounts:         params.Rsns,
		Cutoff:           cutoff,
		PointsGiven:      int(given),
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
func CompetitionWinners(w http.ResponseWriter, r *http.Request) {
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
func CompetitionTeamPosition(w http.ResponseWriter, r *http.Request) {
}
