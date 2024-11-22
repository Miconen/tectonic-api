package handlers

import (
	"tectonic-api/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

var queries *database.Queries

func InitHandlers(conn *pgxpool.Pool) {
  queries = database.New(conn)
}
