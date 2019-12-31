package user

import (
	"context"
	"fmt"

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
		return fmt.Errorf("validater validate: %w", err)
	}

	if err := s.Repository.UniqueUsername(ctx, f.Username); err != nil {
		return fmt.Errorf("unique username: %w", err)
	}

	if err := s.Repository.UniqueEmail(ctx, f.Email); err != nil {
		return fmt.Errorf("unique email: %w", err)
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}

	nu := NewUser{
		Username:     f.Username,
		Email:        f.Email,
		PasswordHash: string(pw),
		RoleID:       f.RoleID,
	}

	if err := s.Repository.Create(ctx, &nu, u); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// Find finds a user by id.
func (s *Service) Find(ctx context.Context, id int) (*User, error) {
	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}

// Update updates a user by id.
func (s *Service) Update(ctx context.Context, id int, f *Form) (*User, error) {
	if err := s.Validater.Validate(ctx, f); err != nil {
		return nil, fmt.Errorf("validate user: %w", err)
	}

	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password hash: %w", err)
	}

	u.Username = f.Username
	u.Email = f.Email
	u.PasswordHash = string(pw)
	u.Role.ID = f.RoleID

	if err := s.Repository.Update(ctx, id, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return u, nil
}

// Delete deletes a user.
func (s *Service) Delete(ctx context.Context, id int) error {
	u, err := s.Repository.Find(ctx, id)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	if err := s.Repository.Delete(ctx, u.ID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

// List shows all users.
func (s *Service) List(ctx context.Context) (*Users, error) {
	var users Users
	if err := s.Repository.List(ctx, &users); err != nil {
		return nil, fmt.Errorf("list of users: %w", err)
	}

	return &users, nil
}
