package handlers

import "net/http"

// @Summary Get RSN related information by guild and user ID
// @Description Get RSN related details by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/rsn [GET]
func GetRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user := User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} User
// @Success 201 {object} Guild
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/rsn [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user := User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Remove RSN from guild and user
// @Description Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 429 {object} Error
// @Failure 500 {object} Error
// @Router /api/v1/rsn [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (interface{}, error) {

		user := User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return user, nil
	}

	httpHandler(w, r, h, p)
}
