package list

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

// Repository allows to work with the database.
type Repository interface {
	ListCategories(ctx context.Context, cat *category.Categories) error
}

// Service is a use case for users showing.
type Service struct {
	Repository
}

// NewService factory prepares service for all futher operations.
func NewService(r Repository) *Service {
	s := Service{
		Repository: r,
	}

	return &s
}

// List shows all categories.
func (s *Service) List(ctx context.Context) (*category.Categories, error) {
	var categories category.Categories
	if err := s.Repository.ListCategories(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "list of categories")
	}

	return &categories, nil
}
