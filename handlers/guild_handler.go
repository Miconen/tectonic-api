package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get a guild by ID
// @Description Get guild details by unique guild Snowflake (ID)
// @Tags Guild
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Success 200 {object} models.Guild
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id} [GET]
func GetGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		fmt.Fprintf(os.Stderr, "No guild id found")
		status = http.StatusBadRequest
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	guild, err := queries.GetGuild(r.Context(), guildId)
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
// @Accept json
// @Produce json
// @Param guild body models.InputGuild true "Guild"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds [POST]
func CreateGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	p := models.InputGuild{
		Multiplier: 1,
	}

	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	_, err = queries.CreateGuild(r.Context(), p.GuildId)
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
// @Router /api/v1/guilds/{guild_id} [DELETE]
func DeleteGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	v := mux.Vars(r)

	guildId, ok := v["guild_id"]
	if !ok {
		fmt.Fprintf(os.Stderr, "No guild id found\n")
		status = http.StatusBadRequest
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	_, err := queries.DeleteGuild(r.Context(), guildId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting guild: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Updates a guild
// @Description Update multiplier and/or time channel for a guild
// @Tags Guild
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param guild body models.UpdateGuild true "Guild"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id} [PUT]
func UpdateGuild(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	v := mux.Vars(r)
	var p database.UpdateGuildParams
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse request body\n")
		status = http.StatusInternalServerError
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	guildId, exists := v["guild_id"]
	if !exists {
		fmt.Fprintf(os.Stderr, "No guild ID found in URL\n")
		status = http.StatusBadRequest
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	if p.GuildID != guildId {
		http.Error(w, fmt.Errorf("guild_id in request body must match URI param").Error(), http.StatusBadRequest)
		return
	}

	_, err = queries.UpdateGuild(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating channel: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
