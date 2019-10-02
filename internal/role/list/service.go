package list

import (
	"context"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
)

// Repository allows to work with the database.
type Repository interface {
	ListRoles(ctx context.Context, roles *role.Roles) error
}

// Service is a use case for users showing.
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

// List shows all roles.
func (s *Service) List(ctx context.Context) (*role.Roles, error) {
	var roles role.Roles
	if err := s.Repository.ListRoles(ctx, &roles); err != nil {
		return nil, errors.Wrap(err, "list of roles")
	}
	return &roles, nil
}
