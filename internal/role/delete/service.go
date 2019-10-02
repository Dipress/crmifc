package delete

import (
	"context"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=delete -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindRole(ctx context.Context, id int) (*role.Role, error)
	DeleteRole(ctx context.Context, id int) error
}

// Service is a use case for role deletion.
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

// Delete deletes a role.
func (s *Service) Delete(ctx context.Context, id int) error {
	rl, err := s.Repository.FindRole(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find role")
	}

	if err := s.Repository.DeleteRole(ctx, rl.ID); err != nil {
		return errors.Wrap(err, "delete role")
	}
	return nil
}
