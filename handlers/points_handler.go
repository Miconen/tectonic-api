package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary		Update a user(s) points
// @Description	Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags			Points
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			user_ids	path		[]string	true	"User IDs"
// @Param			point_event	path		string		true	"Point event"
// @Param			guild		body		models.User	true	"User"
// @Success		200			{object}	models.User
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_ids}/points/{point_event} [PUT]
func UpdatePoints(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	params := database.UpdatePointsByEventParams{
		Event:   p["point_event"],
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	user, err := queries.UpdatePointsByEvent(r.Context(), params)
	ei := database.ClassifyError(err)
	if err != nil {
		handleDatabaseError(*ei, jw, models.ERROR_USER_NOT_FOUND)
		return
	}

	if len(user) == 0 {
		// TODO: check what parameters were wrong, this shouldn't return empty
		jw.WriteError(models.ERROR_TODO)
		return
	}

	jw.WriteResponse(user)
}

// @Summary		Update a user(s) points
// @Description	Update a user(s)' points in our backend by unique user Snowflake (ID)
// @Tags			Points
// @Accept			json
// @Produce		json
// @Param			guild_id	path		string		true	"Guild ID"
// @Param			user_ids	path		[]string	true	"User ID"
// @Param			points		path		string		true	"Points"
// @Param			guild		body		models.User	true	"User"
// @Success		200			{object}	models.Empty
// @Failure		400			{object}	models.Empty
// @Failure		401			{object}	models.Empty
// @Failure		409			{object}	models.Empty
// @Failure		429			{object}	models.Empty
// @Failure		500			{object}	models.Empty
// @Router			/api/v1/guilds/{guild_id}/users/{user_ids}/points/custom/{points} [PUT]
func UpdatePointsCustom(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	params := database.UpdatePointsCustomParams{
		Points:  0,
		UserIds: strings.Split(p["user_ids"], ","),
		GuildID: p["guild_id"],
	}

	points, err := strconv.Atoi(p["points"])
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	params.Points = int32(points)

	rowsaf, err := queries.UpdatePointsCustom(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_TODO)
		return
	}

	if rowsaf == 0 {
		// TODO: check what parameters were wrong, this shouldn't return empty
		jw.WriteError(models.ERROR_TODO)
		return
	}

	jw.WriteResponse(http.NoBody)
}
