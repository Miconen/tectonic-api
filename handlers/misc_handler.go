package handlers

import (
	"context"

	"tectonic-api/database"
	"tectonic-api/models"
)

type (
	GetBossesInput  struct{}
	GetBossesOutput struct {
		Body []database.Boss
	}
)

func (s *Server) GetBosses(ctx context.Context, input *GetBossesInput) (*GetBossesOutput, error) {
	bosses, err := s.queries.GetBosses(ctx)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, models.NewTectonicError(s.getConstraintError(*ei))
	}
	return &GetBossesOutput{Body: bosses}, nil
}

type (
	GetCategoriesInput  struct{}
	GetCategoriesOutput struct {
		Body []database.Category
	}
)

func (s *Server) GetCategories(ctx context.Context, input *GetCategoriesInput) (*GetCategoriesOutput, error) {
	categories, err := s.queries.GetCategories(ctx)
	if ei := database.ClassifyError(err); ei != nil {
		return nil, models.NewTectonicError(s.getConstraintError(*ei))
	}
	return &GetCategoriesOutput{Body: categories}, nil
}
