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
	"github.com/gorilla/mux"
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

func TestFindHandler(t *testing.T) {
	tests := []struct {
		name     string
		findFunc func(ctx context.Context, id int) (*article.Article, error)
		code     int
	}{
		{
			name: "ok",
			findFunc: func(ctx context.Context, id int) (*article.Article, error) {
				return &article.Article{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "internal error",
			findFunc: func(ctx context.Context, id int) (*article.Article, error) {
				return &article.Article{}, errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := FindHandler{findFunc(tc.findFunc)}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type findFunc func(ctx context.Context, id int) (*article.Article, error)

func (f findFunc) Find(ctx context.Context, id int) (*article.Article, error) {
	return f(ctx, id)
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		updateFunc func(ctx context.Context, id int, f *article.Form) (*article.Article, error)
		code       int
	}{
		{
			name: "ok",
			updateFunc: func(ctx context.Context, id int, f *article.Form) (*article.Article, error) {
				return &article.Article{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			updateFunc: func(ctx context.Context, id int, f *article.Form) (*article.Article, error) {
				return &article.Article{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			updateFunc: func(ctx context.Context, id int, f *article.Form) (*article.Article, error) {
				return &article.Article{}, errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := UpdateHandler{updateFunc(tc.updateFunc)}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type updateFunc func(ctx context.Context, id int, f *article.Form) (*article.Article, error)

func (u updateFunc) Update(ctx context.Context, id int, f *article.Form) (*article.Article, error) {
	return u(ctx, id, f)
}
