package role

import (
	"errors"
	"time"
)

// easyjson -all model.go

// ErrNotFound raises when role isn't found in the database.
var ErrNotFound = errors.New("role not found")

// Role constains all role fields.
type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRole contains the information which needs to create a new Role.
type NewRole struct {
	Name string `json:"name"`
}

// Form is a role form.
type Form struct {
	Name string `json:"name"`
}

// Roles contains slice of roles.
type Roles struct {
	Roles []Role `json:"roles"`
}
