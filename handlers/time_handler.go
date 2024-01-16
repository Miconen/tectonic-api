package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids path string true "User IDs"
// @Param time path string true "Time in ticks"
// @Param boss path string true "Boss name"
// @Success 200 {object} models.Empty
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/time [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p, err := utils.ParseParametersURL(r, "time", "boss", "guild_id", "user_ids")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	t, err := strconv.Atoi(p["time"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing (%s) to int: %v\n", p["time"], err)
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	pb, err := database.CheckPb(r.Context(), p)
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

	if pb <= t && pb != 0 {
		status = http.StatusOK
		utils.JsonWriter(pb).IntoHTTP(status)(w, r)
		return
	}

	time, err := database.InsertTime(r.Context(), p)
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
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/time [DELETE]
func RemoveTime(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "run_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	_, err = strconv.Atoi(p["time"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(p["time"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(p["time"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
