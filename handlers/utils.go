package handlers

import (
	"context"
	"tectonic-api/database"
	"tectonic-api/logging"
	"tectonic-api/models"
	"tectonic-api/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbHandler struct {
	database.ErrorInfo
}

func (s *Server) getConstraintError(ei database.ErrorInfo) models.APIV1Error {
	switch ei.Code {
	case "23503": // Foreign constraint violation

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
	case "23505": // Unique constraint violation

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
	return models.ERROR_TODO(400, time.Now().String())
}

// Handlers should delegate database errors to this function always.
// This handler will alwas write the response, so the caller should
// always short circuit the handler.
func (s *Server) handleDatabaseError(ei database.ErrorInfo, jw *utils.JsonWriter) {
	dh := &dbHandler{ErrorInfo: ei}

	if dh.Severity == database.SeverityFatal || dh.Severity == database.SeverityPanic {
		logging.Get().Error("DATABASE FAILURE", "err", dh.Error())
		jw.WriteError(models.ERROR_API_DEAD)
	} else if !dh.Recoverable {
		logging.Get().Error("DATABASE NON-RECOVERABLE ERROR", "err", dh.Error())
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
	} else {
		logging.Get().Warn("QUERY ERROR", "err", dh.Error())
		jw.WriteError(s.getConstraintError(ei))
	}
}

type dbHandlerFunc func(dh *dbHandler, jw *utils.JsonWriter)

func (s *Server) handleDatabaseErrorCustom(ei database.ErrorInfo, jw *utils.JsonWriter, dhc dbHandlerFunc) {
	dh := &dbHandler{ErrorInfo: ei}

	if dh.Severity == database.SeverityFatal || dh.Severity == database.SeverityPanic {
		logging.Get().Error("DATABASE FAILURE", "err", dh.Error())
		jw.WriteError(models.ERROR_API_DEAD)
	} else if !dh.Recoverable {
		logging.Get().Error("DATABASE NON-RECOVERABLE ERROR", "err", dh.Error())
		jw.WriteError(models.ERROR_API_UNAVAILABLE)
	} else {
		logging.Get().Warn("QUERY ERROR", "err", dh.Error())
		dhc(dh, jw)
	}
}

// Checks if guild exists on the database.
func guildExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, guild_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT guild_id FROM guilds WHERE guild_id = $1)", guild_id, models.ERROR_GUILD_NOT_FOUND)
}

// Checks if user exists on the database.
func userExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, user_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT user_id FROM users WHERE user_id = $1)", user_id, models.ERROR_USER_NOT_FOUND)
}

// Checks if event exists on the database.
func eventExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, event_id string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT wom_id FROM event WHERE wom_id = $1)", event_id, models.ERROR_EVENT_NOT_FOUND)
}

// Checks if achievent exists on the database.
func achievementExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, name string) bool {
	return queryExists(ctx, conn, jw, "SELECT EXISTS (SELECT name FROM achievement WHERE name = $1)", name, models.ERROR_ACHIEVEMENT_NOT_FOUND)
}

func queryExists(ctx context.Context, conn *pgxpool.Conn, jw *utils.JsonWriter, sql string, param string, api_err models.APIV1ErrorCode) bool {
	exists := false
	err := conn.QueryRow(ctx, sql, param).Scan(&exists)
	if err != nil {
		jw.WriteError(models.ERROR_API_DEAD)
		return false
	}

	if !exists {
		jw.WriteError(api_err)
		return false
	}

	return true
}
