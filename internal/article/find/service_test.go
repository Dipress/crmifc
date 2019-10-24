package find

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/article"
	"github.com/stretchr/testify/assert"
)

func TestFindArticle(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, id int) (*article.Article, error)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, id int) (*article.Article, error) {
				return &article.Article{}, nil
			},
		},
		{
			name: "repository error",
			repositoryFunc: func(ctx context.Context, id int) (*article.Article, error) {
				return &article.Article{}, errors.New("mock error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := NewService(repositoryFunc(tc.repositoryFunc))

			_, err := s.FindArticle(ctx, 1)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, id int) (*article.Article, error)

func (f repositoryFunc) FindArticle(ctx context.Context, id int) (*article.Article, error) {
	return f(ctx, id)
}
