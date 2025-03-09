package handlers

import (
	"fmt"
	"net/http"
	"os"
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
	status := http.StatusOK

	p := mux.Vars(r)

	id, err := strconv.Atoi(p["competition_id"])
	if err != nil {
		status = http.StatusInternalServerError
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	cutoff, err := strconv.Atoi(p["cutoff"])
	if err != nil {
		status = http.StatusInternalServerError
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	competition, err := utils.GetCompetition(id)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	params := database.GetDetailedUsersByRSNParams{
		GuildID: p["guild_id"],
		Rsns:    make([]string, len(competition.Participations)),
	}

	accounts := make([]string, len(competition.Participations))

	// Get participants valid for points
	for i, val := range competition.Participations {
		accounts[i] = val.Player.DisplayName

		// Skip player if they didn't make the cutoff
		if val.Progress.Gained < float64(cutoff) {
			continue
		}

		params.Rsns[i] = val.Player.DisplayName
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating transaction: %v\n", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
	}

	q := queries.WithTx(tx)
	defer tx.Rollback(r.Context())

	rows, err := q.GetDetailedUsersByRSN(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
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

	points_params := database.UpdatePointsByEventParams{
		Event:   "event_participation",
		GuildID: p["guild_id"],
		UserIds: user_ids,
	}

	points, err := q.UpdatePointsByEvent(r.Context(), points_params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding points to users: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	// Copy updated point values from "points" to "users"
	pointsMap := make(map[string]int)
	for _, user := range points {
		pointsMap[user.UserID] = int(user.Points)
	}

	// We grab a user to track point changes with
	keyUser := points[0].UserID
	var pointsBefore int
	var pointsAfter int

	for i := range users {
		if val, exists := pointsMap[points[i].UserID]; exists {
			// Handle tracking of changing point values to calculate given points
			if users[i].UserID == keyUser {
				pointsBefore = int(users[i].Points)
				pointsAfter = val
			}

			// Replace old point value with new one
			users[i].Points = int32(val)
		}
	}

	// Calculate given points dynamically
	given := pointsAfter - pointsBefore

	tx.Commit(r.Context())

	// TODO: Add points given to response
	response := CompetitionResponse{
		Title:            competition.Title,
		ParticipantCount: competition.ParticipantCount,
		Participants:     users,
		Accounts:         accounts,
		Cutoff:           cutoff,
		PointsGiven:      given,
	}

	utils.JsonWriter(response).IntoHTTP(status)(w, r)
}
