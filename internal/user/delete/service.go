package delete

import (
	"context"

	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=delete -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindUser(ctx context.Context, id int) (*user.User, error)
	DeleteUser(ctx context.Context, id int) error
}

// Service is a use case for user deletion.
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

// Delete deletes a user.
func (s *Service) Delete(ctx context.Context, id int) error {
	u, err := s.Repository.FindUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find user")
	}

	if err := s.Repository.DeleteUser(ctx, u.ID); err != nil {
		return errors.Wrap(err, "delete user")
	}

	return nil
}
