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

type InputPointsEvent struct {
	GuildID string   `json:"guild_id"`
	UserIDs []string `json:"user_ids"`
	Event   string   `json:"event"`
}

func (i InputPointsEvent) GetUserIDs() []string {
	return i.UserIDs
}

func (i InputPointsEvent) GetGuildID() string {
	return i.GuildID
}

type InputPointsCustom struct {
	GuildID string   `json:"guild_id"`
	UserIDs []string `json:"user_ids"`
	Points  int      `json:"points"`
}

func (i InputPointsCustom) GetUserIDs() []string {
	return i.UserIDs
}

func (i InputPointsCustom) GetGuildID() string {
	return i.GuildID
}

// @Summary Update a user(s) points
// @Description Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags Points
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param point_event path string true "Point event"
// @Param guild body models.User true "User"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/points/{point_event} [PUT]
func UpdatePoints(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	v := mux.Vars(r)
	p := database.UpdatePointsByEventParams{}

	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		fmt.Println("Error parsing request body:", err)
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	if p.GuildID != v["guild_id"] {
		http.Error(w, fmt.Errorf("guild_id in request body must match URI param").Error(), http.StatusBadRequest)
		return
	}

	if p.Event != v["point_event"] {
		fmt.Println(p.Event, v["points"])
		http.Error(w, "event in request body must match URI param", http.StatusBadRequest)
		return
	}

	_, err = queries.UpdatePointsByEvent(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating points: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Update a user(s) points
// @Description Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags Points
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param points path string true "Points"
// @Param guild body models.User true "User"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/points/custom/{points} [PUT]
func UpdatePointsCustom(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	v := mux.Vars(r)
	p := database.UpdatePointsCustomParams{}

	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		fmt.Println("Error parsing request body:", err)
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	if p.GuildID != v["guild_id"] {
		http.Error(w, fmt.Errorf("guild_id in request body must match URI param").Error(), http.StatusBadRequest)
		return
	}

	vp, err := strconv.Atoi(v["points"])
	if err != nil {
		status = http.StatusBadRequest
		fmt.Println("Error parsing points:", err)
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	if p.Points != int32(vp) {
		http.Error(w, fmt.Errorf("points in request body must match URI param").Error(), http.StatusBadRequest)
		return
	}

	_, err = queries.UpdatePointsCustom(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating points: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
