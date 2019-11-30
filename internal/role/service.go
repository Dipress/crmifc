package role

import (
	"context"

	"github.com/pkg/errors"
)

// go:generate mockgen -source=service.go -package=role -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	Create(ctx context.Context, f *NewRole, rol *Role) error
	Find(ctx context.Context, id int) (*Role, error)
	Update(ctx context.Context, id int, rl *Role) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, roles *Roles) error
}

// Validater validates role fields.
type Validater interface {
	Validate(ctx context.Context, form *Form) error
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
func (s *Service) Create(ctx context.Context, f *Form) (*Role, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	var nr NewRole
	nr.Name = f.Name

	var rol Role
	if err := s.Repository.Create(ctx, &nr, &rol); err != nil {
		return nil, errors.Wrap(err, "repository create role")
	}
	return &rol, nil
}

// Find finds a role.
func (s *Service) Find(ctx context.Context, id int) (*Role, error) {
	r, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}
	return r, nil
}

// Update updates a role.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*Role, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validater validate")
	}

	rl, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "repository find role")
	}

	rl.Name = f.Name

	if err := s.Repository.Update(ctx, id, rl); err != nil {
		return nil, errors.Wrap(err, "update role")
	}
	return rl, nil
}

// Delete deletes a role.
func (s *Service) Delete(ctx context.Context, id int) error {
	rl, err := s.Repository.Find(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find role")
	}

	if err := s.Repository.Delete(ctx, rl.ID); err != nil {
		return errors.Wrap(err, "delete role")
	}
	return nil
}

// List shows all roles.
func (s *Service) List(ctx context.Context) (*Roles, error) {
	var roles Roles
	if err := s.Repository.List(ctx, &roles); err != nil {
		return nil, errors.Wrap(err, "list of roles")
	}
	return &roles, nil
}
