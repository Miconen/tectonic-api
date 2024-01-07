package handlers

import (
	"net/http"
	"tectonic-api/models"
)

// @Summary Get RSN related information by guild and user ID
// @Description Get RSN related details by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/rsn [GET]
func GetRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return models.Body{Content: user}, http.StatusOK, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/rsn [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return models.Body{Content: user}, http.StatusCreated, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Remove RSN from guild and user
// @Description Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 429 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/rsn [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
		"rsn":      r.URL.Query().Get("rsn"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		user := models.User{
			UserId:  p["user_id"],
			GuildId: p["guild_id"],
			Points:  789,
		}

		return models.Body{Content: user}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}
