package update

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=update -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindCategory(ctx context.Context, id int) (*category.Category, error)
	UpdateCategory(ctx context.Context, id int, cat *category.Category) error
}

// Validater validates category fields.
type Validater interface {
	Validate(ctx context.Context, form *category.Form) error
}

// Service is a use case for category validation and updation.
type Service struct {
	Repository
	Validater
}

// NewService factory prepares service for all futher operations.
func NewService(r Repository, v Validater) *Service {
	s := Service{
		Repository: r,
		Validater:  v,
	}

	return &s
}

// Update updates a category.
func (s *Service) Update(ctx context.Context, id int, f *category.Form) (*category.Category, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	cat, err := s.Repository.FindCategory(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find category")
	}

	cat.Name = f.Name

	if err := s.Repository.UpdateCategory(ctx, id, cat); err != nil {
		return nil, errors.Wrap(err, "update category")
	}

	return cat, nil
}
