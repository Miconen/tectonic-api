package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/database"
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

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn body models.InputRSN true "RSN"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/wom/competition/{competition_id}/cutoff/{cutoff} [GET]
func EndCompetition(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	id, err := strconv.Atoi(p["competition_id"])
	if err != nil {
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(err)
		return
	}

	cutoff, err := strconv.Atoi(p["cutoff"])
	if err != nil {
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(err)
		return
	}

	competition, err := utils.GetCompetition(id)
	if err != nil {
		log.Error("Error fetching WOM competition", "error", err)
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	if len(competition.Participations) == 0 {
		log.Info("no participations found")
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
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
		log.Info("no participations found with specified cutoff")
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(err)
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		log.Error("Error creating transaction", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	rows, err := q.GetDetailedUsersByRSN(r.Context(), params)
	if err != nil {
		log.Error("Error fetching user", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
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
	if err != nil {
		log.Error("Error fetching points value to give by event", "error", err)
		jw.SetStatus(http.StatusInternalServerError)
		jw.WriteResponse(http.NoBody)
		return
	}

	points_params := database.UpdatePointsByEventParams{
		Event:   "event_participation",
		GuildID: p["guild_id"],
		UserIds: user_ids,
	}

	points, err := q.UpdatePointsByEvent(r.Context(), points_params)
	if err != nil {
		log.Error("Error adding points to users", "error", err)
		jw.SetStatus(http.StatusNotFound)
		jw.WriteResponse(http.NoBody)
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
