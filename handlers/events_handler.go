package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"tectonic-api/database"
	"tectonic-api/models"
	"tectonic-api/utils"

	"github.com/gorilla/mux"
)

// @Summary		Get the guild's events
// @Description	Get the events that the guild have created
// @Tags			Event
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	database.Event
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/events [GET]
func (s *Server) GetEvents(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(s.queries.GetGuildEvents, r.Context(), p["guild_id"])
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	jw.WriteResponse(events)
}

// @Summary		Get the guild's event
// @Description	Get the event that the guild have created
// @Tags			Event
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	models.DetailedEvent
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/events/{event_id} [GET]
func (s *Server) GetDetailedEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(s.queries.GetEventParticipation, r.Context(), p["event_id"])
	if err != nil {
		s.handleDatabaseError(*err, jw)
		return
	}

	res := &models.DetailedEvent{
		Participations: utils.MapField(events, func(p database.GetEventParticipationRow) models.EventParticipation {
			return models.EventParticipation{
				UserId:    p.UserID,
				Placement: int(p.Placement),
			}
		}),
	}

	jw.WriteResponse(res)
}

// @Summary		Register a guild event
// @Description	Register a guild event present in the WOM APi
// @Tags			Event
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		201			{object}	database.Event
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/events [POST]
func (s *Server) RegisterEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	p := mux.Vars(r)
	body := models.InputEvent{
		PositionCutoff: 3,
	}

	if err := utils.ParseAndValidateRequestBody(w, r, &body); err != nil {
		return
	}

	c, err := s.womClient.GetCompetition(body.EventId)
	if err != nil {
		var apiErr *utils.WomAPIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case http.StatusNotFound:
				jw.WriteError(models.ERROR_WOMID_NOT_FOUND)
			case http.StatusGatewayTimeout:
				jw.WriteError(models.ERROR_WOM_UNAVAILABLE)
			default:
				jw.WriteError(models.ERROR_API_UNAVAILABLE)
			}
		} else {
			// network or JSON decoding error
			jw.WriteError(models.ERROR_API_UNAVAILABLE)
		}
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}
	defer tx.Rollback(r.Context())

	q := s.queries.WithTx(tx)
	ei := database.WrapExec(q.CreateEvent, r.Context(), database.CreateEventParams{
		Name:           c.Title,
		WomID:          fmt.Sprintf("%d", c.ID),
		GuildID:        p["guild_id"],
		PositionCutoff: 0,
	})
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	if c.Type == "classic" {
		ei = database.WrapExec(q.InsertEventParticipants, r.Context(), database.InsertEventParticipantsParams{
			ParticipantIds: utils.MapField(c.Participations, func(p models.Participations) string {
				return fmt.Sprintf("%d", p.PlayerID)
			}),
			GuildID: p["guild_id"],
			WomID:   fmt.Sprintf("%d", c.ID),
		})
		if ei != nil {
			s.handleDatabaseError(*ei, jw)
			return
		}
	} else if c.Type == "team" && len(body.TeamNames) != 0 {
		ids := make([]string, len(c.Participations))
		positions := make([]int32, len(c.Participations))

		pos_map := make(map[string]int32)
		for i := range body.TeamNames {
			pos_map[body.TeamNames[i]] = int32(i)
		}

		for i := range c.Participations {
			ids[i] = fmt.Sprintf("%d", c.Participations[i].PlayerID)
			positions[i] = pos_map[c.Participations[i].TeamName]
		}

		ei = database.WrapExec(q.InsertEventTeams, r.Context(), database.InsertEventTeamsParams{
			ParticipantIds:        ids,
			ParticipantPlacements: positions,
			GuildID:               p["guild_id"],
			WomID:                 fmt.Sprintf("%d", c.ID),
		})
		if ei != nil {
			s.handleDatabaseError(*ei, jw)
			return
		}
	} else {
		jw.WriteError(models.ERROR_WRONG_BODY)
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}

	jw.WriteResponse(http.NoBody)
}

// @Summary		Delete a guild event
// @Description	Delete a guild event registered in the API
// @Tags			Event
// @Produce		json
// @Param			guild_id	path		string	true	"Guild ID"
// @Success		200			{object}	database.Event
// @Failure		400			{object}	models.ErrorResponse
// @Failure		401			{object}	models.ErrorResponse
// @Failure		404			{object}	models.ErrorResponse
// @Failure		429			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/api/v1/guilds/{guild_id}/events/{event_id} [DELETE]
func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	ei := database.WrapExec(s.queries.DeleteEvent, r.Context(), p["event_id"])
	if ei != nil {
		s.handleDatabaseError(*ei, jw)
		return
	}

	jw.WriteResponse(http.NoBody)
}
