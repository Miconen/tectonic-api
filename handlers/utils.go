package handlers

import (
	"errors"
	"net/http"

	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
)

func (s *Server) getConstraintError(ei database.ErrorInfo) models.APIV1Error {
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
		case "user_achievement":
			return models.ERROR_USER_ACHIEVEMENT_EXISTS
		case "users":
			return models.ERROR_USER_EXISTS
		}
	}

	logging.Get().Error("Database untreated error", "error_info", ei)
	return models.ERROR_API_UNAVAILABLE
}

func (s *Server) dbError(ei database.ErrorInfo) *models.TectonicError {
	if ei.Severity == database.SeverityFatal || ei.Severity == database.SeverityPanic {
		logging.Get().Error("DATABASE FAILURE", "err", ei.Error())
		return models.NewTectonicError(models.ERROR_API_DEAD)
	}
	if !ei.Recoverable {
		logging.Get().Error("DATABASE NON-RECOVERABLE ERROR", "err", ei.Error())
		return models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
	}
	logging.Get().Warn("QUERY ERROR", "err", ei.Error())
	return models.NewTectonicError(s.getConstraintError(ei))
}

func (s *Server) womError(err error) *models.TectonicError {
	var apiErr *utils.WomAPIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return models.NewTectonicError(models.ERROR_WOMID_NOT_FOUND)
		case http.StatusGatewayTimeout:
			return models.NewTectonicError(models.ERROR_WOM_UNAVAILABLE)
		}
	}
	return models.NewTectonicError(models.ERROR_API_UNAVAILABLE)
}
