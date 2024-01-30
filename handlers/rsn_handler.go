package handlers

import (
	"fmt"
	"net/http"
	"os"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get RSN related information by guild and user ID
// @Description Get RSN related details by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string false "Guild ID"
// @Param user_id path string false "User ID"
// @Success 200 {object} []models.RSN
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id}/rsns [GET]
func GetRSNs(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	rsns, err := database.SelectRsns(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting RSN: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(rsns).IntoHTTP(status)(w, r)
}

// @Summary Link an RSN to a user
// @Description Link an RSN to a guild and user in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn body models.InputRSN true "RSN"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id}/rsns [POST]
func CreateRSN(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	v := mux.Vars(r)
	p := models.InputRSN{}
	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}
	if p.GuildId != v["guild_id"] {
		http.Error(w, fmt.Errorf("guild_id in request body must match URI param").Error(), http.StatusBadRequest)
		return
	}

	wid, err := utils.GetWomId(p.RSN)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.InsertRsn(r.Context(), p, wid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating RSN: %v\n", err)
		// TODO: Handle 404 Not Found errors
		status = http.StatusConflict
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Remove RSN from guild and user
// @Description Delete a RSN in our backend by unique guild and user Snowflake (ID)
// @Tags RSN
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Param rsn path string true "RSN"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id}/rsns/{rsn} [DELETE]
func RemoveRSN(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p := mux.Vars(r)

	err := database.DeleteRsn(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting RSN: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
