package category

import (
	"context"

	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=category -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	Create(ctx context.Context, f *NewCategory, cat *Category) error
	Find(ctx context.Context, id int) (*Category, error)
	Update(ctx context.Context, id int, cat *Category) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, cat *Categories) error
}

// Validater validates role fields.
type Validater interface {
	Validate(ctx context.Context, form *Form) error
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
func (s *Service) Create(ctx context.Context, f *Form) (*Category, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	var nc NewCategory
	nc.Name = f.Name

	var cat Category
	if err := s.Repository.Create(ctx, &nc, &cat); err != nil {
		return nil, errors.Wrap(err, "repository create category")
	}

	return &cat, nil
}

// Find finds a category.
func (s *Service) Find(ctx context.Context, id int) (*Category, error) {
	c, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return c, nil
}

// Update updates a category.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*Category, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	cat, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find category")
	}

	cat.Name = f.Name

	if err := s.Repository.Update(ctx, id, cat); err != nil {
		return nil, errors.Wrap(err, "update category")
	}

	return cat, nil
}

// Delete deletes a category.
func (s *Service) Delete(ctx context.Context, id int) error {
	cat, err := s.Repository.Find(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find category")
	}

	if err := s.Repository.Delete(ctx, cat.ID); err != nil {
		return errors.Wrap(err, "delete category")
	}

	return nil
}

// List shows all categories.
func (s *Service) List(ctx context.Context) (*Categories, error) {
	var categories Categories
	if err := s.Repository.List(ctx, &categories); err != nil {
		return nil, errors.Wrap(err, "list of categories")
	}

	return &categories, nil
}
