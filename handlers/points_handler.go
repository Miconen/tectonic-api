package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Update a user(s) points
// @Description Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags Points
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param point_event path string true "Point event"
// @Param guild body models.User true "User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_ids}/points/{point_event} [PUT]
func UpdatePoints(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	params := database.UpdatePointsByEventParams{
		Event:   p["point_event"],
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	user, err := queries.UpdatePointsByEvent(r.Context(), params)
	if err != nil {
		log.Error("Error updating points", "error", err)
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(user)
}

// @Summary Update a user(s) points
// @Description Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags Points
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param points path string true "Points"
// @Param guild body models.User true "User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_ids}/points/custom/{points} [PUT]
func UpdatePointsCustom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	params := database.UpdatePointsCustomParams{
		Points:  0,
		UserIds: strings.Split(p["user_ids"], ","),
		GuildID: p["guild_id"],
	}

	points, err := strconv.Atoi(p["points"])
	if err != nil {
		log.Error("Error parsing points", "error", err)
		jw.SetStatus(http.StatusBadRequest)
		jw.WriteResponse(err)
		return
	}

	params.Points = int32(points)

	user, err := queries.UpdatePointsCustom(r.Context(), params)
	if err != nil {
		log.Error("Error updating points", "error", err)
		jw.SetStatus(http.StatusNotFound)
	}

	jw.WriteResponse(user)
}
