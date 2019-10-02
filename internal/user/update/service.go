package update

import (
	"context"

	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// go:generate mockgen -source=service.go -package=update -destination=service.mock.go

// Repository allows working with a database.
type Repository interface {
	FindUser(ctx context.Context, id int) (*user.User, error)
	UpdateUser(ctx context.Context, id int, u *user.User) error
}

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, form *user.Form) error
}

// Service is a use case for user updation.
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

// Update updates a user by id.
func (s *Service) Update(ctx context.Context, id int, f *user.Form) error {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return errors.Wrap(err, "validate user")
	}

	u, err := s.Repository.FindUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find user")
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	u.Username = f.Username
	u.Email = f.Email
	u.PasswordHash = string(pw)
	u.Role.ID = f.RoleID

	if err := s.Repository.UpdateUser(ctx, id, u); err != nil {
		return errors.Wrap(err, "update user")
	}

	return nil
}
