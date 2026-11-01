package handlers

import (
	"context"
	"errors"
	"net/http"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
)

func (s *Server) getConstraintError(ctx context.Context, ei database.ErrorInfo) models.APIV1Error {
	switch ei.Code {
	case "23503":
		c := s.constraintsMap[ei.Err.ConstraintName]
		switch c.ForeignTable {
		case "achievement":
			return models.ERROR_ACHIEVEMENT_NOT_FOUND
		case "bosses":
			return models.ERROR_BOSS_NOT_FOUND
		case "categories":
			return models.ERROR_CATEGORY_NOT_FOUND
		case "event":
			return models.ERROR_EVENT_NOT_FOUND
		case "event_participant":
			return models.ERROR_PARTICIPATION_NOT_FOUND
		case "guild_bosses":
			return models.ERROR_GUILD_BOSS_NOT_FOUND
		case "guild_categories":
			return models.ERROR_GUILD_CATEGORY_NOT_FOUND
		case "guilds":
			return models.ERROR_GUILD_NOT_FOUND
		case "point_sources":
			return models.ERROR_POINT_SOURCE_NOT_FOUND
		case "rsn":
			return models.ERROR_RSN_NOT_FOUND
		case "teams":
			return models.ERROR_TEAM_NOT_FOUND
		case "times":
			return models.ERROR_TIME_NOT_FOUND
		case "records":
			return models.ERROR_RECORD_NOT_FOUND
		case "guild_ranks":
			return models.ERROR_GUILD_RANK_NOT_FOUND
		case "user_achievement":
			return models.ERROR_USER_ACHIEVEMENT_NOT_FOUND
		case "users":
			return models.ERROR_USER_NOT_FOUND
		}
	case "23505":
		c := s.constraintsMap[ei.Err.ConstraintName]
		switch c.Table {
		case "achievement":
			return models.ERROR_ACHIEVEMENT_EXISTS
		case "bosses":
			return models.ERROR_BOSS_EXISTS
		case "categories":
			return models.ERROR_CATEGORY_EXISTS
		case "event":
			return models.ERROR_EVENT_EXISTS
		case "event_participant":
			return models.ERROR_PARTICIPATION_EXISTS
		case "guild_bosses":
			return models.ERROR_GUILD_BOSS_EXISTS
		case "guild_categories":
			return models.ERROR_GUILD_CATEGORY_EXISTS
		case "guilds":
			return models.ERROR_GUILD_EXISTS
		case "point_sources":
			return models.ERROR_POINT_SOURCE_EXISTS
		case "rsn":
			return models.ERROR_RSN_EXISTS
		case "teams":
			return models.ERROR_TEAM_EXISTS
		case "times":
			return models.ERROR_TIME_EXISTS
		case "records":
			return models.ERROR_RECORD_EXISTS
		case "guild_ranks":
			return models.ERROR_GUILD_RANK_EXISTS
		case "user_achievement":
			return models.ERROR_USER_ACHIEVEMENT_EXISTS
		case "users":
			return models.ERROR_USER_EXISTS
		}
	}

	logging.From(ctx).Error("Untreated database error", "err", ei.Error(), "code", ei.Code)
	return models.ERROR_API_UNAVAILABLE
}

// dbError maps a database error to an API error and logs it with the
// request's ID, so a failure can be found from the caller's Error ID.
func (s *Server) dbError(ctx context.Context, ei database.ErrorInfo) *models.TectonicError {
	log := logging.From(ctx)
	if ei.Severity == database.SeverityFatal || ei.Severity == database.SeverityPanic {
		log.Error("Database failure", "err", ei.Error())
		return models.NewTectonicError(models.ERROR_API_DEAD)
	}
	if !ei.Recoverable {
		log.Error("Non-recoverable database error", "err", ei.Error())
		return models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	code := s.getConstraintError(ctx, ei)
	log.Warn("Query error", "err", ei.Error(), "code", ei.Code, "constraint", constraintName(ei))
	return models.NewTectonicError(code)
}

func constraintName(ei database.ErrorInfo) string {
	if ei.Err == nil {
		return ""
	}
	return ei.Err.ConstraintName
}

func (s *Server) womError(err error, notFound models.APIV1Error) *models.TectonicError {
	var apiErr *utils.WomAPIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return models.NewTectonicError(notFound)
		case http.StatusGatewayTimeout:
			return models.NewTectonicError(models.ERROR_WOM_RATE_LIMITED)
		case http.StatusInternalServerError:
			return models.NewTectonicError(models.ERROR_WOM_UNAVAILABLE)
		}
	}
	return models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
}
