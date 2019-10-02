package find

import (
	"context"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
)

// Repository allows to work with the database.
type Repository interface {
	FindRole(ctx context.Context, id int) (*role.Role, error)
}

// Service is a use case for role finding.
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

// Find finds a role.
func (s *Service) Find(ctx context.Context, id int) (*role.Role, error) {
	r, err := s.Repository.FindRole(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}
	return r, nil
}
