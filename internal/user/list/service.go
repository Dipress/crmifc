package list

import (
	"context"

	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
)

// Repository allows to work with the database.
type Repository interface {
	ListUsers(ctx context.Context, users *user.Users) error
}

// Service is a use case for users showing.
type Service struct {
	Repository
}

// NewService factory prepares service for all further operations.
func NewService(r Repository) *Service {
	s := Service{
		Repository: r,
	}

	return &s
}

// List shows all users.
func (s *Service) List(ctx context.Context) (*user.Users, error) {
	var users user.Users
	if err := s.Repository.ListUsers(ctx, &users); err != nil {
		return nil, errors.Wrap(err, "list of users")
	}

	return &users, nil
}
