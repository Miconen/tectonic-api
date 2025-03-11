package handlers

import (
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var queries *database.Queries
var log = utils.NewLogger()

func InitHandlers(conn *pgxpool.Pool) {
  queries = database.New(conn)
}
