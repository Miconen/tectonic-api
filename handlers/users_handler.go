package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get one or more users by ID(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids path string true "User ID(s)"
// @Success 200 {object} database.User[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_ids} [GET]
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	params := database.GetUsersByIdParams{
		GuildID: p["guild_id"],
		UserIDs: strings.Split(p["user_ids"], ","),
	}

	user, err := queries.GetUsersById(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(user).IntoHTTP(status)(w, r)
}

// @Summary Get one or more users by RSN(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param rsns path string true "User RSN(s)"
// @Success 200 {object} database.User[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/rsn/{rsns} [GET]
func GetUsersByRsn(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	params := database.GetUsersByRsnParams{
		GuildID: p["guild_id"],
		RSNs: strings.Split(p["rsns"], ","),
	}

	user, err := queries.GetUsersByRsn(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(user).IntoHTTP(status)(w, r)
}

// @Summary Get one or more users by WomID(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param wom_id path string true "User WomID(s)"
// @Success 200 {object} database.User[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/wom/{wom_ids} [GET]
func GetUsersByWom(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	params := database.GetUsersByWomParams{
		GuildID: p["guild_id"],
		WomIDs: strings.Split(p["wom_ids"], ","),
	}

	user, err := queries.GetUsersByWom(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	utils.JsonWriter(user).IntoHTTP(status)(w, r)
}

// @Summary Create / Initialize a new user
// @Description Initialize a user in our backend by unique user Snowflake (ID)
// @Tags User
// @Accept json
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param guild body models.InputUser true "User"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	v := mux.Vars(r)
	p := models.InputUser{}
	err := utils.ParseRequestBody(w, r, &p)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}
	p.GuildId = v["guild_id"]

	wid, err := utils.GetWomId(p.RSN)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	err = database.InsertUser(r.Context(), p, wid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting user: %v\n", err)
		if err.Error() == database.ERROR_UNACTIVATED_GUILD {
			status = http.StatusNotFound
		} else {
			status = http.StatusConflict
		}
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Delete a user from guild
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_id path string true "User ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id} [DELETE]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p := mux.Vars(r)

	err := database.DeleteUser(r.Context(), p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
