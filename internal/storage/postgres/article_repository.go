package postgres

import (
	"context"
	"database/sql"

	"github.com/dipress/crmifc/internal/article"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ArticleRepository holds CRUD actions for article.
type ArticleRepository struct {
	db *sqlx.DB
}

//NewArticleRepository factory prepares the repository to work.
func NewArticleRepository(db *sql.DB) *ArticleRepository {
	r := ArticleRepository{
		db: sqlx.NewDb(db, driverName),
	}

	return &r
}

const createArticleQuery = `INSERT INTO 
	articles (user_id, category_id, title, body) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, category_id, title, body, created_at, updated_at`

// Create inserts a new category into the database.
func (r *ArticleRepository) Create(ctx context.Context, f *article.NewArticle, art *article.Article) error {
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

// Find finds a article by id.
func (r *ArticleRepository) Find(ctx context.Context, id int) (*article.Article, error) {
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

// Update updates article by id.
func (r *ArticleRepository) Update(ctx context.Context, id int, a *article.Article) error {
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

// Delete deletes article by id.
func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
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

const listArticleQuery = `SELECT * FROM articles`

// List shows all articles.
func (r *ArticleRepository) List(ctx context.Context, articles *article.Articles) error {
	rows, err := r.db.QueryxContext(ctx, listArticleQuery)
	if err != nil {
		return errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	for rows.Next() {
		var a article.Article
		if err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Title,
			&a.Body,
			&a.CategoryID,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return errors.Wrap(err, "articles query row scan on loop")
		}

		articles.Articles = append(articles.Articles, a)
	}

	return nil
}
