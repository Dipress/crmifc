package delete

import (
	"context"

	"github.com/dipress/crmifc/internal/article"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=delete -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindArticle(ctx context.Context, id int) (*article.Article, error)
	DeleteArticle(ctx context.Context, id int) error
}

// Service is a use case for category deletion.
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
	art, err := s.Repository.FindArticle(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find article")
	}

	if err := s.Repository.DeleteArticle(ctx, art.ID); err != nil {
		return errors.Wrap(err, "delete category")
	}

	return nil
}
