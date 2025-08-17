package handlers

import (
	"context"
	"tectonic-api/config"
	"tectonic-api/database"
	"tectonic-api/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	queries        *database.Queries
	pool           *pgxpool.Pool
	womClient      *utils.WomClient
	constraintsMap map[string]database.ConstraintDetail
	config         *config.Config
}

func NewServer(pool *pgxpool.Pool, wom *utils.WomClient, cfg *config.Config) (*Server, error) {
	queries := database.New(pool)

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	constraintsMap, err := database.GetConstraintsTable(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	return &Server{
		pool:           pool,
		queries:        queries,
		womClient:      wom,
		constraintsMap: constraintsMap,
		config:         cfg,
	}, nil
}

func (s *Server) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *Server) Queries() *database.Queries {
	return s.queries
}

func (s *Server) WomClient() *utils.WomClient {
	return s.womClient
}

func (s *Server) ConstraintsMap() map[string]database.ConstraintDetail {
	return s.constraintsMap
}

func (s *Server) Config() *config.Config {
	return s.config
}
