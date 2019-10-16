package create

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=create -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	CreateCategory(ctx context.Context, f *category.NewCategory, cat *category.Category) error
}

// Validater validates role fields.
type Validater interface {
	Validate(ctx context.Context, form *category.Form) error
}

// Service is a use case for category creation.
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

// Create creates a category.
func (s *Service) Create(ctx context.Context, f *category.Form) (*category.Category, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	var nc category.NewCategory
	nc.Name = f.Name

	var cat category.Category
	if err := s.Repository.CreateCategory(ctx, &nc, &cat); err != nil {
		return nil, errors.Wrap(err, "repository create category")
	}

	return &cat, nil
}
