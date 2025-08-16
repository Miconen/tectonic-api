package handlers

import (
	"context"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var queries *database.Queries
var womClient *utils.WomClient

func InitHandlers(conn *pgxpool.Pool, wom *utils.WomClient) {
	c, err := conn.Acquire(context.Background())
	if err != nil {
		panic("failed to get conn from pool while initializing handlers")
	}
	defer c.Release()

	InitDatabaseHandler(context.Background(), c)
	queries = database.New(conn)
	womClient = wom
}
