package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param time body models.InputTime true "Time"
// @Success 200 {object} models.Empty
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guilds/{guild_id}/times [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	v := mux.Vars(r)

	p := models.InputTime{}
	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	pb, err := database.CheckPb(r.Context(), v["guild_id"], p)
	if err != nil {
		if pb == -1 {
			fmt.Fprintf(os.Stderr, "Error fetching pb: %v\n", err)
			status = http.StatusInternalServerError
			utils.JsonWriter(err).IntoHTTP(status)(w, r)
			return
		}
		if pb == 0 {
			fmt.Fprintf(os.Stderr, "Pb is null: %v\n", err)
		}
	}

	if pb <= p.Time && pb != 0 {
		status = http.StatusOK
		utils.JsonWriter(pb).IntoHTTP(status)(w, r)
		return
	}

	time, err := database.InsertTime(r.Context(), v["guild_id"], p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting time: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(time).IntoHTTP(status)(w, r)
}

// @Summary Remove time from guilds best times
// @Description Delete a time in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param time_id path string true "Time ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guilds/{guild_id}/times/{time_id} [DELETE]
func RemoveTime(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	v := mux.Vars(r)

	_, err := strconv.Atoi(v["time_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteTime(r.Context(), v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting time: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
