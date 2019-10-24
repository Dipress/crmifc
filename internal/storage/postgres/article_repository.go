package postgres

import (
	"context"
	"database/sql"

	"github.com/dipress/crmifc/internal/article"
	"github.com/pkg/errors"
)

const createArticleQuery = `INSERT INTO 
	articles (user_id, category_id,  title, body) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, category_id, title, body, created_at, updated_at`

// CreateArticle inserts a new category into the database.
func (r *Repository) CreateArticle(ctx context.Context, f *article.NewArticle, art *article.Article) error {
	if err := r.db.QueryRowContext(ctx, createArticleQuery, f.UserID, f.CategoryID, f.Title, f.Body).
		Scan(
			&art.ID,
			&art.UserID,
			&art.CategoryID,
			&art.Title,
			&art.Body,
			&art.CreatedAt,
			&art.UpdatedAt,
		); err != nil {
		return errors.Wrap(err, "query context scan")
	}

	return nil
}

const findArticleQuery = `SELECT id, user_id, category_id, title, body, created_at, updated_at FROM articles WHERE id = $1`

// FindArticle finds a article by id.
func (r *Repository) FindArticle(ctx context.Context, id int) (*article.Article, error) {
	var a article.Article
	if err := r.db.QueryRowContext(ctx, findArticleQuery, id).
		Scan(
			&a.ID,
			&a.UserID,
			&a.CategoryID,
			&a.Title,
			&a.Body,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
		if err == sql.ErrNoRows {
			return nil, article.ErrNotFound
		}

		return nil, errors.Wrap(err, "query row scan")
	}

	return &a, nil
}
