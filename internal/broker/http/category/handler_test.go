package category

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/validation"
)

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		createFunc func(ctx context.Context, f *category.Form) (*category.Category, error)
		code       int
	}{
		{
			name: "ok",
			createFunc: func(ctx context.Context, f *category.Form) (*category.Category, error) {
				return &category.Category{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation",
			createFunc: func(ctx context.Context, f *category.Form) (*category.Category, error) {
				return &category.Category{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internl error",
			createFunc: func(ctx context.Context, f *category.Form) (*category.Category, error) {
				return &category.Category{}, errors.New("mock error")
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

			r := httptest.NewRequest("POST", "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type createFunc func(ctx context.Context, f *category.Form) (*category.Category, error)

func (c createFunc) Create(ctx context.Context, f *category.Form) (*category.Category, error) {
	return c(ctx, f)
}
