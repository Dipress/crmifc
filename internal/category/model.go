package category

import (
	"errors"
	"time"
)

// easyjson -all model.go

var (
	// ErrNotFound raises when category isn't found in the database.
	ErrNotFound = errors.New("category not found")
	// ErrNameExists returns when given name is already
	// present in database.
	ErrNameExists = errors.New("name already exists")
)

// Category contains all user field.
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Form is a category form.
type Form struct {
	Name string `json:"name"`
}

// NewCategory contains the information which needs to create a new Category.
type NewCategory struct {
	Name string `json:"name"`
}

// Categories contains slice of categories.
type Categories struct {
	Categories []Category `json:"categories"`
}
