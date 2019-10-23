package article

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/validation"
)

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		createFunc func(ctx context.Context, f *article.Form) (*article.Article, error)
		code       int
	}{
		{
			name: "ok",
			createFunc: func(ctx context.Context, f *article.Form) (*article.Article, error) {
				return &article.Article{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			createFunc: func(ctx context.Context, f *article.Form) (*article.Article, error) {
				return &article.Article{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			createFunc: func(ctx context.Context, f *article.Form) (*article.Article, error) {
				return &article.Article{}, errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := CreateHandler{createFunc(tc.createFunc)}
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type createFunc func(ctx context.Context, f *article.Form) (*article.Article, error)

func (c createFunc) Create(ctx context.Context, f *article.Form) (*article.Article, error) {
	return c(ctx, f)
}
