package find

import (
	"context"

	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
)

// Repository allows working with a database.
type Repository interface {
	FindUser(ctx context.Context, id int) (*user.User, error)
}

// Service is a use case for user finding.
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

// Find finds a user by id.
func (s *Service) Find(ctx context.Context, id int) (*user.User, error) {
	u, err := s.Repository.FindUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find user")
	}
	return u, nil
}
