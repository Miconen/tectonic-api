package handlers

import (
	"net/http"
	"strconv"
)

// @Summary Add a new best time to guild
// @Description Add a new time to a guild in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Time
// @Success 201 {object} Time
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/time [POST]
func CreateTime(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"time":      r.URL.Query().Get("time"),
		"boss_name": r.URL.Query().Get("boss_name"),
		"run_id":    r.URL.Query().Get("run_id"),
		"date":      r.URL.Query().Get("date"),
		"team":      r.URL.Query().Get("team"),
	}

	h := func(r *http.Request) (interface{}, error) {
		timeInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		runIdInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		dateInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		team := Users{}
		for i := 0; i < 10; i++ {
			user := User{
				UserId:  "Hello",
				GuildId: "World",
				Points:  789,
			}

			team.Users = append(team.Users, user)
		}

		time := Time{
			Time:     timeInt,
			BossName: p["boss_name"],
			RunId:    runIdInt,
			Date:     dateInt,
			Team:     team,
		}

		return time, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Remove time from guilds best times
// @Description Delete a time in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Time
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/time [DELETE]
func RemoveTime(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"time":      r.URL.Query().Get("time"),
		"boss_name": r.URL.Query().Get("boss_name"),
		"run_id":    r.URL.Query().Get("run_id"),
		"date":      r.URL.Query().Get("date"),
		"team":      r.URL.Query().Get("team"),
	}

	h := func(r *http.Request) (interface{}, error) {
		timeInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		runIdInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		dateInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		team := Users{}
		for i := 0; i < 10; i++ {
			user := User{
				UserId:  "Hello",
				GuildId: "World",
				Points:  789,
			}

			team.Users = append(team.Users, user)
		}

		time := Time{
			Time:     timeInt,
			BossName: p["boss_name"],
			RunId:    runIdInt,
			Date:     dateInt,
			Team:     team,
		}

		return time, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update times channel from guild
// @Description Update where time related embeds are located in our backend by unique guild Snowflake (ID)
// @Tags Time
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} Time
// @Success 201 {object} Time
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /v1/time [PUT]
func UpdateTimesChannel(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"time":      r.URL.Query().Get("time"),
		"boss_name": r.URL.Query().Get("boss_name"),
		"run_id":    r.URL.Query().Get("run_id"),
		"date":      r.URL.Query().Get("date"),
		"team":      r.URL.Query().Get("team"),
	}

	h := func(r *http.Request) (interface{}, error) {
		timeInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		runIdInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		dateInt, err := strconv.Atoi(p["time"])
		if err != nil {
			return nil, err
		}

		team := Users{}
		for i := 0; i < 10; i++ {
			user := User{
				UserId:  "Hello",
				GuildId: "World",
				Points:  789,
			}

			team.Users = append(team.Users, user)
		}

		time := Time{
			Time:     timeInt,
			BossName: p["boss_name"],
			RunId:    runIdInt,
			Date:     dateInt,
			Team:     team,
		}

		return time, nil
	}

	httpHandler(w, r, h, p)
}
