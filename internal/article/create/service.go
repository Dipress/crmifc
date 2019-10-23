package create

import (
	"context"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=create -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindByEmail(ctx context.Context, email string, usr *user.User) error
	CreateArticle(ctx context.Context, f *article.NewArticle, art *article.Article) error
}

// Validater validates article fields.
type Validater interface {
	Validate(ctx context.Context, form *article.Form) error
}

// Service is a use case for article creation.
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
func (s *Service) Create(ctx context.Context, f *article.Form) (*article.Article, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater error")
	}

	claims, _ := auth.FromContext(ctx)

	var u user.User
	if err := s.Repository.FindByEmail(ctx, claims.Subject, &u); err != nil {
		return nil, errors.Wrap(err, "repository find user by email")
	}

	na := article.NewArticle{
		CategoryID: f.CategoryID,
		Title:      f.Title,
		Body:       f.Body,
		UserID:     u.ID,
	}

	var a article.Article
	if err := s.Repository.CreateArticle(ctx, &na, &a); err != nil {
		return nil, errors.Wrap(err, "repository create article")
	}
	return &a, nil
}
