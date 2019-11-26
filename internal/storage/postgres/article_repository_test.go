package postgres

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/article"
)

func TestCreateArticle(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewArticleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.Log("\ttest:0\tshould create the article into the database")
		{
			na := article.NewArticle{
				UserID:     1,
				CategoryID: 1,
				Title:      "article title",
				Body:       "article body",
			}

			var art article.Article
			err := r.Create(ctx, &na, &art)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if art.ID == 0 {
				t.Error("expected to parse returned id")
			}

		}
	}
}

func TestFindArticle(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r := NewArticleRepository(db)

		na := article.NewArticle{
			UserID:     2,
			CategoryID: 3,
			Title:      "my new title",
			Body:       "my new body",
		}

		var art article.Article
		err := r.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould find the article into the database")
		{
			_, err := r.Find(ctx, art.ID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestUpdateArticle(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewArticleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		na := article.NewArticle{
			UserID:     2,
			CategoryID: 12,
			Title:      "my new title",
			Body:       "my new body",
		}

		var art article.Article
		err := r.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould update the article into the database")
		{
			art.Title = "my update title"
			art.Body = "my update body"
			art.CategoryID = 13

			err := r.Update(ctx, 1, &art)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}

func TestArticleDelete(t *testing.T) {
	t.Log("with initialized repository")
	{
		db, teardown := postgresDB(t)
		defer teardown()

		r := NewArticleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		na := article.NewArticle{
			UserID:     5,
			CategoryID: 32,
			Title:      "my new title",
			Body:       "my new body",
		}

		var art article.Article
		err := r.Create(ctx, &na, &art)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		t.Log("\ttest:0\tshould delete the article into the database")
		{
			err := r.Delete(ctx, art.ID)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}
}
