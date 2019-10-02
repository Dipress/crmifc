package user

import (
	"errors"
	"time"

	"github.com/dipress/crmifc/internal/role"
)

// easyjson -all model.go

var (
	// ErrNotFound raises when user isn't found in the database.
	ErrNotFound = errors.New("user not found")
	// ErrUsernameExists returns when given username is already
	// present in database.
	ErrUsernameExists = errors.New("username already exists")
	// ErrEmailExists returns when given email is already
	// present in database.
	ErrEmailExists = errors.New("email already exists")
)

// User contains all user field.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Role         role.Role `json:"role"`
}

// Form is a user form.
type Form struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

// NewUser contains the information which needs to create a new User.
type NewUser struct {
	RoleID       int
	Username     string
	Email        string
	PasswordHash string
}

// Users contains slice of users.
type Users struct {
	Users []User `json:"users"`
}
