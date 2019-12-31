package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dipress/crmifc/internal/category"
	"github.com/jmoiron/sqlx"
)

// CategoryRepository holds CRUD actions.
type CategoryRepository struct {
	db *sqlx.DB
}

//NewCategoryRepository factory prepares the repository to work.
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	r := CategoryRepository{
		db: sqlx.NewDb(db, driverName),
	}

	return &r
}

const createCategoryQuery = `INSERT INTO 
	categories (name) 
	VALUES ($1) 
	RETURNING id, name, created_at, updated_at`

// Create inserts a new category into the database.
func (r *CategoryRepository) Create(ctx context.Context, f *category.NewCategory, cat *category.Category) error {
	if err := r.db.QueryRowContext(ctx, createCategoryQuery, f.Name).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
		return fmt.Errorf("query context scan: %w", err)
	}

	return nil
}

const findCategoryQuery = `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`

// Find finds a category by id.
func (r *CategoryRepository) Find(ctx context.Context, id int) (*category.Category, error) {
	var cat category.Category

	if err := r.db.QueryRowContext(ctx, findCategoryQuery, id).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, category.ErrNotFound
		}
		return nil, fmt.Errorf("query row scan: %w", err)
	}

	return &cat, nil
}

const updateCategoryQuery = `UPDATE categories SET name=:name, updated_at=now() WHERE id=:id`

// Update updates a category by id.
func (r *CategoryRepository) Update(ctx context.Context, id int, cat *category.Category) error {
	stmt, err := r.db.PrepareNamed(updateCategoryQuery)
	if err != nil {
		return fmt.Errorf("prepare named: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id":   id,
		"name": cat.Name,
	}); err != nil {
		if err == sql.ErrNoRows {
			return category.ErrNotFound
		}

		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

const deleteCategoryQuery = `DELETE FROM categories WHERE id=:id`

// Delete deletes category by id.
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareNamed(deleteCategoryQuery)
	if err != nil {
		return fmt.Errorf("prepare named: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		if err == sql.ErrNoRows {
			return category.ErrNotFound
		}

		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

const listCategoryQuery = `SELECT * FROM categories`

// List shows all categories.
func (r *CategoryRepository) List(ctx context.Context, cat *category.Categories) error {
	rows, err := r.db.QueryxContext(ctx, listCategoryQuery)
	if err != nil {
		return fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c category.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return fmt.Errorf("categories query row scan on loop: %w", err)
		}

		cat.Categories = append(cat.Categories, c)
	}

	return nil
}
