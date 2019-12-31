package category

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("validater validate: %w", err)
	}

	var nc NewCategory
	nc.Name = f.Name

	var cat Category
	if err := s.Repository.Create(ctx, &nc, &cat); err != nil {
		return nil, fmt.Errorf("repository create category: %w", err)
	}

	return &cat, nil
}

// Find finds a category.
func (s *Service) Find(ctx context.Context, id int) (*Category, error) {
	c, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository find: %w", err)
	}

	return c, nil
}

// Update updates a category.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*Category, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, fmt.Errorf("validater validate: %w", err)
	}

	cat, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository find category: %w", err)
	}

	cat.Name = f.Name

	if err := s.Repository.Update(ctx, id, cat); err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return cat, nil
}

// Delete deletes a category.
func (s *Service) Delete(ctx context.Context, id int) error {
	cat, err := s.Repository.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("find category: %w", err)
	}

	if err := s.Repository.Delete(ctx, cat.ID); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}

// List shows all categories.
func (s *Service) List(ctx context.Context) (*Categories, error) {
	var categories Categories
	if err := s.Repository.List(ctx, &categories); err != nil {
		return nil, fmt.Errorf("list of categories: %w", err)
	}

	return &categories, nil
}
