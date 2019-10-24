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

const updateArticleQuery = `UPDATE articles SET user_id=:user_id, category_id=:category_id, title=:title, body=:body, updated_at=now() WHERE id=:id`

// UpdateArticle updates article by id.
func (r *Repository) UpdateArticle(ctx context.Context, id int, a *article.Article) error {
	stmt, err := r.db.PrepareNamed(updateArticleQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id":          id,
		"user_id":     a.UserID,
		"category_id": a.CategoryID,
		"title":       a.Title,
		"body":        a.Body,
	}); err != nil {
		if err == sql.ErrNoRows {
			return article.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}
	return nil
}

const deleteArticleQuery = `DELETE FROM articles WHERE id=:id`

// DeleteArticle deletes article by id.
func (r *Repository) DeleteArticle(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareNamed(deleteArticleQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		if err == sql.ErrNoRows {
			return article.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}

	return nil
}
