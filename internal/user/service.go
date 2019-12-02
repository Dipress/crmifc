package user

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// go:generate mockgen -source=service.go -package=user -destination=service.mock.go

// Repository allows to work with the database.
type Repository interface {
	Create(ctx context.Context, f *NewUser, usr *User) error
	UniqueUsername(ctx context.Context, username string) error
	UniqueEmail(ctx context.Context, email string) error
	Find(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, id int, u *User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, users *Users) error
}

// Validater validates user fields.
type Validater interface {
	Validate(ctx context.Context, form *Form) error
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
func (s *Service) Create(ctx context.Context, f *Form, u *User) error {
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

	nu := NewUser{
		Username:     f.Username,
		Email:        f.Email,
		PasswordHash: string(pw),
		RoleID:       f.RoleID,
	}

	if err := s.Repository.Create(ctx, &nu, u); err != nil {
		return errors.Wrap(err, "create user")
	}

	return nil
}

// Find finds a user by id.
func (s *Service) Find(ctx context.Context, id int) (*User, error) {
	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find user")
	}
	return u, nil
}

// Update updates a user by id.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*User, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, errors.Wrap(err, "validate user")
	}

	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "find user")
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u.Username = f.Username
	u.Email = f.Email
	u.PasswordHash = string(pw)
	u.Role.ID = f.RoleID

	if err := s.Repository.Update(ctx, id, u); err != nil {
		return nil, errors.Wrap(err, "update user")
	}

	return u, nil
}

// Delete deletes a user.
func (s *Service) Delete(ctx context.Context, id int) error {
	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return errors.Wrap(err, "find user")
	}

	if err := s.Repository.Delete(ctx, u.ID); err != nil {
		return errors.Wrap(err, "delete user")
	}

	return nil
}

// List shows all users.
func (s *Service) List(ctx context.Context) (*Users, error) {
	var users Users
	if err := s.Repository.List(ctx, &users); err != nil {
		return nil, errors.Wrap(err, "list of users")
	}

	return &users, nil
}
