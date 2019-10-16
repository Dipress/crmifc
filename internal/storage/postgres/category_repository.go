package postgres

import (
	"context"

	"github.com/dipress/crmifc/internal/category"
	"github.com/pkg/errors"
)

const createCategoryQuery = `INSERT INTO 
	categories (name) 
	VALUES ($1) 
	RETURNING id, name, created_at, updated_at`

// CreateCategory inserts a new category into the database.
func (r *Repository) CreateCategory(ctx context.Context, f *category.NewCategory, cat *category.Category) error {
	if err := r.db.QueryRowContext(ctx, createCategoryQuery, f.Name).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
		return errors.Wrap(err, "query context scan")
	}

	return nil
}
