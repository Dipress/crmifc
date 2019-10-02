package create

import (
	"context"

	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// go:generate mockgen -source=service.go -package=create -destination=service.mock.go

// Repository allows working with a database.
type Repository interface {
	CreateUser(ctx context.Context, f *user.NewUser, usr *user.User) error
	UniqueUsername(ctx context.Context, username string) error
	UniqueEmail(ctx context.Context, email string) error
}

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, form *user.Form) error
}

// Service is a use case for user creation.
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

// Create creates a user.
func (s *Service) Create(ctx context.Context, f *user.Form, u *user.User) error {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return errors.Wrap(err, "validater validate")
	}

	if err := s.Repository.UniqueUsername(ctx, f.Username); err != nil {
		return errors.Wrap(err, "unique username")
	}

	if err := s.Repository.UniqueEmail(ctx, f.Email); err != nil {
		return errors.Wrap(err, "unique email")
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	nu := user.NewUser{
		Username:     f.Username,
		Email:        f.Email,
		PasswordHash: string(pw),
		RoleID:       f.RoleID,
	}

	if err := s.Repository.CreateUser(ctx, &nu, u); err != nil {
		return errors.Wrap(err, "create user")
	}

	return nil
}
