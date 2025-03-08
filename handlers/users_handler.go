package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary Get one or more users by ID(s)
// @Description Get user details by unique user Snowflake (ID)
// @Tags Users
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param user_ids path string true "User ID(s)"
// @Success 200 {object} database.DetailedUSer[]
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_ids} [GET]
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	p := mux.Vars(r)

	params := database.GetDetailedUsersParams{
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	rows, err := queries.GetDetailedUsers(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching user: %v\n", err)
		status = http.StatusNotFound
		utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
		return
	}

	users := database.NewDetailedUserFromRows(rows)
	utils.JsonWriter(users).IntoHTTP(status)(w, r)
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
		Rsns:    strings.Split(p["rsns"], ","),
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
		WomIds:  strings.Split(p["wom_ids"], ","),
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
// @Param user_id path string true "User ID"
// @Param rsn body string true "RSN"
// @Success 201 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 409 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/{user_id} [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusCreated

	v := mux.Vars(r)
	params := database.CreateUserParams{
		GuildID: v["guild_id"],
		UserID:  v["user_id"],
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Error decoding request body: ", err)
		status = http.StatusInternalServerError
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	wid, err := utils.GetWomId(params.Rsn)
	if err != nil {
		status = http.StatusBadRequest
		utils.JsonWriter(err).IntoHTTP(status)(w, r)
		return
	}

	params.WomID = wid

	user, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting user: %v\n", err)
		if err.Error() == database.ERROR_UNACTIVATED_GUILD {
			status = http.StatusNotFound
		} else {
			status = http.StatusConflict
		}
	}

	utils.JsonWriter(user).IntoHTTP(status)(w, r)
}

// @Summary Delete a user from guild by User ID
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
func RemoveUserById(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p := mux.Vars(r)

	params := database.DeleteUserByIdParams{
		GuildID: p["guild_id"],
		UserID:  p["user_id"],
	}

	err := queries.DeleteUserById(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Delete a user from guild by RSN
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param rsn path string true "RSN"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/rsn/{rsn} [DELETE]
func RemoveUserByRsn(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p := mux.Vars(r)

	params := database.DeleteUserByRsnParams{
		GuildID: p["guild_id"],
		Rsn:     p["rsn"],
	}

	err := queries.DeleteUserByRsn(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}

// @Summary Delete a user from guild by Wom ID
// @Description Delete a user in our backend by unique user and guild Snowflake (ID)
// @Tags User
// @Produce json
// @Param guild_id path string true "Guild ID"
// @Param wom_id path string true "Wom ID"
// @Success 204 {object} models.Empty
// @Failure 400 {object} models.Empty
// @Failure 401 {object} models.Empty
// @Failure 404 {object} models.Empty
// @Failure 429 {object} models.Empty
// @Failure 500 {object} models.Empty
// @Router /api/v1/guilds/{guild_id}/users/wom/{wom_id} [DELETE]
func RemoveUserByWom(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNoContent

	p := mux.Vars(r)

	params := database.DeleteUserByWomParams{
		GuildID: p["guild_id"],
		WomID:   p["wom_id"],
	}

	err := queries.DeleteUserByWom(r.Context(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user: %v\n", err)
		status = http.StatusNotFound
	}

	utils.JsonWriter(http.NoBody).IntoHTTP(status)(w, r)
}
