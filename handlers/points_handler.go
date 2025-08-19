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
func (s *Server) UpdatePoints(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)
	params := database.UpdatePointsByEventParams{
		Event:   p["point_event"],
		GuildID: p["guild_id"],
		UserIds: strings.Split(p["user_ids"], ","),
	}

	user, err := s.queries.UpdatePointsByEvent(r.Context(), params)
	ei := database.ClassifyError(err)
	if err != nil {
		s.handleDatabaseErrorCustom(*ei, jw, func(dh *dbHandler, jw *utils.JsonWriter) {
			switch dh.Code {
			case "23502":
				jw.WriteError(models.ERROR_POINT_SOURCE_NOT_FOUND)
			default:
				jw.WriteError(s.getConstraintError(*ei))
			}
		})
		return
	}

	if len(user) == 0 {
		jw.WriteError(models.ERROR_POINT_SOURCE_NOT_FOUND)
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
func (s *Server) UpdatePointsCustom(w http.ResponseWriter, r *http.Request) {
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

	user, err := s.queries.UpdatePointsCustom(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if len(user) == 0 {
		jw.WriteError(models.ERROR_POINT_SOURCE_NOT_FOUND)
		return
	}

	jw.WriteResponse(user)
}

// @Summary		Get the guild's point sources
// @Description	Get the point sources that the guild has created
// @Tags			Points
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	database.Event
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/points [GET]
func (s *Server) GetPointSources(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(s.queries.GetGuildPointSources, r.Context(), p["guild_id"])
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(events)
}

// @Summary		Update a guild point source
// @Description	Update a guilds points source
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
// @Router			/api/v1/guilds/{guild_id}/points/{point_source}/{points} [PUT]
func (s *Server) UpdateGuildPointSource(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)

	p := mux.Vars(r)

	params := database.UpdateGuildPointSourceParams{
		Points:      0,
		GuildID:     p["guild_id"],
		PointSource: p["point_source"],
	}

	points, err := strconv.Atoi(p["points"])
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_PARAMS)
		return
	}

	params.Points = int32(points)

	rowsaf, err := s.queries.UpdateGuildPointSource(r.Context(), params)
	ei := database.ClassifyError(err)
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if rowsaf == 0 {
		jw.WriteError(models.ERROR_POINT_SOURCE_NOT_FOUND)
		return
	}

	jw.WriteResponse(http.NoBody)
}
