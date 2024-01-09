package handlers

import (
	"net/http"
	"strconv"
	"tectonic-api/models"
)

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 201 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 409 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/time [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"time":      r.URL.Query().Get("time"),
		"boss_name": r.URL.Query().Get("boss_name"),
		"run_id":    r.URL.Query().Get("run_id"),
		"date":      r.URL.Query().Get("date"),
		"team":      r.URL.Query().Get("team"),
	}

	h := func(r *http.Request) (models.Body, int, error) {
		timeInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err
		}

		runIdInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err

		}

		dateInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err
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

		return models.Body{Content: time}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Remove time from guilds best times
// @Description Delete a time in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/time [DELETE]
func RemoveTime(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"time":      r.URL.Query().Get("time"),
		"boss_name": r.URL.Query().Get("boss_name"),
		"run_id":    r.URL.Query().Get("run_id"),
		"date":      r.URL.Query().Get("date"),
		"team":      r.URL.Query().Get("team"),
	}

	h := func(r *http.Request) (models.Body, int, error) {
		timeInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err
		}

		runIdInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err
		}

		dateInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return models.Body{}, http.StatusBadRequest, err
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

		return models.Body{Content: time}, http.StatusNoContent, nil

	}

	httpHandler(w, r, h, p)
}
