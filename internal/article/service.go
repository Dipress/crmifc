package article

import (
	"context"

	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=article -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	Create(ctx context.Context, f *NewArticle, art *Article) error
	Find(ctx context.Context, id int) (*Article, error)
	Update(ctx context.Context, id int, a *Article) error
	Delete(ctx context.Context, id int) error
}

// Validater validates article fields.
type Validater interface {
	Validate(ctx context.Context, form *Form) error
}

// Service is a use case for all actions of the article.
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

// Create creates a new article.
func (s *Service) Create(ctx context.Context, f *Form) (*Article, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater error")
	}

	claims, _ := auth.FromContext(ctx)

	na := NewArticle{
		CategoryID: f.CategoryID,
		Title:      f.Title,
		Body:       f.Body,
		UserID:     claims.User.ID,
	}

	var a Article
	if err := s.Repository.Create(ctx, &na, &a); err != nil {
		return nil, errors.Wrap(err, "repository create article")
	}
	return &a, nil
}

// Find finds a article by id.
func (s *Service) Find(ctx context.Context, id int) (*Article, error) {
	a, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find article")
	}
	return a, nil
}

// Update updates a article.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*Article, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	claims, _ := auth.FromContext(ctx)

	a, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find article")
	}

	a.UserID = claims.User.ID
	a.CategoryID = f.CategoryID
	a.Title = f.Title
	a.Body = f.Body

	if err := s.Repository.Update(ctx, id, a); err != nil {
		return nil, errors.Wrap(err, "update article")
	}

	return a, nil
}

// Delete deletes a article.
func (s *Service) Delete(ctx context.Context, id int) error {
	art, err := s.Repository.Find(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find article")
	}

	if err := s.Repository.Delete(ctx, art.ID); err != nil {
		return errors.Wrap(err, "delete category")
	}

	return nil
}
