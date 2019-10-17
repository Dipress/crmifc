package postgres

import (
	"context"
	"database/sql"

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

const findCategoryQuery = `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`

// Find Category finds a category by id
func (r *Repository) FindCategory(ctx context.Context, id int) (*category.Category, error) {
	var cat category.Category

	if err := r.db.QueryRowContext(ctx, findCategoryQuery, id).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, category.ErrNotFound
		}
		return nil, errors.Wrap(err, "query row scan")
	}

	return &cat, nil
}
