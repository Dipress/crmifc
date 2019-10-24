package update

import (
	"context"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=update -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindByEmail(ctx context.Context, email string, usr *user.User) error
	FindArticle(ctx context.Context, id int) (*article.Article, error)
	UpdateArticle(ctx context.Context, id int, a *article.Article) error
}

// Validater validates article fields.
type Validater interface {
	Validate(ctx context.Context, form *article.Form) error
}

// Service is a use case for article updation.
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

// Update updates a article.
func (s *Service) Update(ctx context.Context, id int, f *article.Form) (*article.Article, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	claims, _ := auth.FromContext(ctx)

	var u user.User
	if err := s.Repository.FindByEmail(ctx, claims.Subject, &u); err != nil {
		return nil, errors.Wrap(err, "repository find user by email")
	}

	a, err := s.Repository.FindArticle(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find article")
	}

	a.UserID = u.ID
	a.CategoryID = f.CategoryID
	a.Title = f.Title
	a.Body = f.Body

	if err := s.Repository.UpdateArticle(ctx, id, a); err != nil {
		return nil, errors.Wrap(err, "update article")
	}

	return a, nil
}
