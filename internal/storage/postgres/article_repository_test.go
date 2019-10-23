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

		r := NewRepository(db)

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
			err := r.CreateArticle(ctx, &na, &art)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if art.ID == 0 {
				t.Error("expected to parse returned id")
			}

		}
	}

}
