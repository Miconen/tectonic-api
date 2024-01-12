package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/models"
	"tectonic-api/utils"
)

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 403 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/time [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p, err := utils.ParseParametersURL(r, "time", "boss_name", "run_id", "date", "team")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	timeInt, err := strconv.Atoi(p["time"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	runIdInt, err := strconv.Atoi(p["date"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dateInt, err := strconv.Atoi(p["time"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	team := models.Users{}
	for i := 0; i < 10; i++ {
		user := models.User{
			UserId:  "Hello",
			GuildId: "World",
			Points:  789,
		}

		team.Users = append(team.Users, user)
	}

	time := models.Time{
		Time:     timeInt,
		BossName: p["boss_name"],
		RunId:    runIdInt,
		Date:     dateInt,
		Team:     team,
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
// @Failure 403 {object} models.Empty
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
