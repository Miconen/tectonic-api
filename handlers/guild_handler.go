package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/utils"
)

// @Summary Get a guild by ID
// @Description Get guild details by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id query string false "Guild ID"
// @Success 200 {object} models.Guild
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guild [GET]
func GetGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p, err := utils.ParseParametersURL(r, "guild_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	guild, err := database.SelectGuild(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting guild: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	// Write JSON response
	utils.JsonWriter(guild).IntoHTTP(status)(w, r)
}

// @Summary Create / Initialize a guild
// @Description Initialize a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guild [POST]
func CreateGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p, err := utils.ParseParametersURL(r, "guild_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.InsertGuild(r.Context(), p["guild_id"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating guild: %v\n", err)
		status = http.StatusConflict
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Delete a guild
// @Description Delete a guild in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guild [DELETE]
func RemoveGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.DeleteGuild(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting guild: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Update times channel from guild
// @Description Update where time related embeds are located in our backend by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guild/times [PUT]
func UpdateTimesChannel(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id", "pb_channel_id")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	g, f, err := utils.ExtractByClone(p, "guild_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = database.UpdateGuild(r.Context(), g, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating channel: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Update multiplier for guild
// @Description Update guild point multiplier by guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /v1/guild/multiplier [PUT]
func UpdateMultiplier(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p, err := utils.ParseParametersURL(r, "guild_id", "multiplier")
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	g, f, err := utils.ExtractByClone(p, "guild_id")
	if err != nil {
		status = http.StatusInternalServerError
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.UpdateGuild(r.Context(), g, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating multiplier: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
