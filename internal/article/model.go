package article

import (
	"errors"
	"time"
)

// easyjson -all model.go

var (
	// ErrNotFound raises when article not found in the database.
	ErrNotFound = errors.New("article not found")
)

// Article contains all article field.
type Article struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CategoryID int       `json:"categort_id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NewArticle contains the information which needs to create a new Article.
type NewArticle struct {
	UserID     int
	CategoryID int
	Title      string
	Body       string
}

// Form is a article form.
type Form struct {
	UserID     int    `json:"user"`
	CategoryID int    `json:"category_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

// Articles contains slice of posts.
type Articles struct {
	Articles []Article `json:"articles"`
}
