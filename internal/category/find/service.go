package find

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

// Repository allows to work with the database.
type Repository interface {
	FindCategory(ctx context.Context, id int) (*category.Category, error)
}

// Service is a use case for category finding.
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

// Find finds a category.
func (s *Service) Find(ctx context.Context, id int) (*category.Category, error) {
	c, err := s.Repository.FindCategory(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return c, nil
}
