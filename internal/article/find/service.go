package find

import (
	"context"

	"github.com/dipress/crmifc/internal/article"
	"github.com/pkg/errors"
)

// Repository allows working with a database.
type Repository interface {
	FindArticle(ctx context.Context, id int) (*article.Article, error)
}

// Service is a use case for article finding.
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

// Find finds a article by id.
func (s *Service) Find(ctx context.Context, id int) (*article.Article, error) {
	a, err := s.Repository.FindArticle(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find article")
	}
	return a, nil
}
