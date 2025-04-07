package handlers

import (
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
func GetEvents(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(queries.GetGuildEvents, r.Context(), p["guild_id"])
	if err != nil {
		handleDatabaseError(*err, jw, models.ERROR_API_DEAD)
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
func GetDetailedEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	events, err := database.WrapQuery(queries.GetEventParticipation, r.Context(), p["event_id"])
	if err != nil {
		handleDatabaseError(*err, jw, models.ERROR_API_DEAD)
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
// @Router			/api/v1/guilds/{guild_id}/events [PUT]
func RegisterEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusCreated)

	p := mux.Vars(r)
	b := models.InputEvent{
		PositionCutoff: 3,
	}

	err := utils.ParseRequestBody(w, r, &b)
	if err != nil {
		jw.WriteError(models.ERROR_WRONG_BODY)
		return
	}

	c, err := utils.GetCompetition(b.EventId)
	if err != nil {
		jw.WriteError(models.ERROR_WOM_UNAVAILABLE)
		return
	}

	tx, err := database.CreateTx(r.Context())
	if err != nil {
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
		return
	}
	defer tx.Rollback(r.Context())

	q := queries.WithTx(tx)
	ei := database.WrapExec(q.CreateEvent, r.Context(), database.CreateEventParams{
		Name:           c.Title,
		WomID:          fmt.Sprintf("%d", c.ID),
		GuildID:        p["guild_id"],
		PositionCutoff: 0,
	})
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_API_DEAD)
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
			handleDatabaseError(*ei, jw, models.ERROR_API_DEAD)
			return
		}
	} else if c.Type == "team" && len(b.TeamNames) != 0 {
		ids := make([]string, len(c.Participations))
		positions := make([]int32, len(c.Participations))

		pos_map := make(map[string]int32)
		for i := range b.TeamNames {
			pos_map[b.TeamNames[i]] = int32(i)
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
			handleDatabaseError(*ei, jw, models.ERROR_API_DEAD)
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
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	jw := utils.NewJsonWriter(w, r, http.StatusOK)
	p := mux.Vars(r)

	ei := database.WrapExec(queries.DeleteEvent, r.Context(), p["event_id"])
	if ei != nil {
		handleDatabaseError(*ei, jw, models.ERROR_API_DEAD)
		return
	}

	jw.WriteResponse(http.NoBody)
}
