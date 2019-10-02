package create

import (
	"context"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=create -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	CreateRole(ctx context.Context, f *role.NewRole, rol *role.Role) error
}

// Validater validates role fields.
type Validater interface {
	Validate(ctx context.Context, form *role.Form) error
}

// Service is a use case for role creation.
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

// Create creates a role.
func (s *Service) Create(ctx context.Context, f *role.Form) (*role.Role, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	var nr role.NewRole
	nr.Name = f.Name

	var rol role.Role
	if err := s.Repository.CreateRole(ctx, &nr, &rol); err != nil {
		return nil, errors.Wrap(err, "repository create role")
	}
	return &rol, nil
}
