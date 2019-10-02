package update

import (
	"context"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=update -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	FindRole(ctx context.Context, id int) (*role.Role, error)
	UpdateRole(ctx context.Context, id int, rl *role.Role) error
}

// Validater validates role fields.
type Validater interface {
	Validate(ctx context.Context, form *role.Form) error
}

// Service is a use case for role validation and updation.
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

// Update updates a role.
func (s *Service) Update(ctx context.Context, id int, f *role.Form) (*role.Role, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	rl, err := s.Repository.FindRole(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find role")
	}

	rl.Name = f.Name

	if err := s.Repository.UpdateRole(ctx, id, rl); err != nil {
		return nil, errors.Wrap(err, "update role")
	}
	return rl, nil
}
