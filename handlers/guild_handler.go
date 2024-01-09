package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
)

// @Summary Get a guild by ID
// @Description Get guild details by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/guild [GET]
func GetGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		guild := models.Guild{
			GuildId: p["guild_id"],
		}
		guild, err := database.SelectGuild(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting guild: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{Content: guild}, http.StatusOK, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Create / Initialize a guild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 201 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 409 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/guild [POST]
func CreateGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		guild := models.Guild{
			GuildId: p["guild_id"],
		}
		err := database.InsertGuild(guild)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating guild: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusCreated, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Delete a guild
// @Description Delete a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/guild [DELETE]
func RemoveGuild(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id": r.URL.Query().Get("guild_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		err := database.DeleteGuild(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting guild: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update times channel from guild
// @Description Update where time related embeds are located in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/guild/times [PUT]
func UpdateTimesChannel(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id":      r.URL.Query().Get("guild_id"),
		"pb_channel_id": r.URL.Query().Get("pb_channel_id"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		c := map[string]interface{}{
			"pb_channel_id": p["pb_channel_id"],
		}

		err := database.UpdateGuild(p["guild_id"], c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating channel: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}

// @Summary Update multiplier for guild
// @Description Update guild point multiplier by guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Body
// @Failure 400 {object} models.Body
// @Failure 403 {object} models.Body
// @Failure 404 {object} models.Body
// @Failure 429 {object} models.Body
// @Failure 500 {object} models.Body
// @Router /v1/guild/multiplier [PUT]
func UpdateMultiplier(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"guild_id":   r.URL.Query().Get("guild_id"),
		"multiplier": r.URL.Query().Get("multiplier"),
	}

	h := func(r *http.Request) (models.Body, int, error) {

		c := map[string]interface{}{
			"multiplier": p["multiplier"],
		}

		err := database.UpdateGuild(p["guild_id"], c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating multiplier: %v\n", err)
			return models.Body{}, http.StatusInternalServerError, err
		}

		return models.Body{}, http.StatusNoContent, nil
	}

	httpHandler(w, r, h, p)
}
