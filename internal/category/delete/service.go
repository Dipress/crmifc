package delete

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=delete -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindCategory(ctx context.Context, id int) (*category.Category, error)
	DeleteCategory(ctx context.Context, id int) error
}

// Service is a use case for category validation and updation.
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

// Delete deletes a category.
func (s *Service) Delete(ctx context.Context, id int) error {
	cat, err := s.Repository.FindCategory(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find category")
	}

	if err := s.Repository.DeleteCategory(ctx, cat.ID); err != nil {
		return errors.Wrap(err, "delete category")
	}

	return nil
}
