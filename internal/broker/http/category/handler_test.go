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
	"github.com/gorilla/mux"
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

			r := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("{}"))

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

func TestFindHandler(t *testing.T) {
	tests := []struct {
		name     string
		findFunc func(ctx context.Context, id int) (*category.Category, error)
		code     int
	}{
		{
			name: "ok",
			findFunc: func(ctx context.Context, id int) (*category.Category, error) {
				return &category.Category{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "internal error",
			findFunc: func(ctx context.Context, id int) (*category.Category, error) {
				return &category.Category{}, errors.New("mock error")
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

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		updateFunc func(ctx context.Context, id int, f *category.Form) (*category.Category, error)
		code       int
	}{
		{
			name: "ok",
			updateFunc: func(ctx context.Context, id int, f *category.Form) (*category.Category, error) {
				return &category.Category{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			updateFunc: func(ctx context.Context, id int, f *category.Form) (*category.Category, error) {
				return &category.Category{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			updateFunc: func(ctx context.Context, id int, f *category.Form) (*category.Category, error) {
				return &category.Category{}, errors.New("mock error")
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

type updateFunc func(ctx context.Context, id int, f *category.Form) (*category.Category, error)

func (u updateFunc) Update(ctx context.Context, id int, f *category.Form) (*category.Category, error) {
	return u(ctx, id, f)
}

type findFunc func(ctx context.Context, id int) (*category.Category, error)

func (f findFunc) Find(ctx context.Context, id int) (*category.Category, error) {
	return f(ctx, id)
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		deleteFunc func(ctx context.Context, id int) error
		code       int
	}{
		{
			name: "ok",
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "repository error",
			deleteFunc: func(ctx context.Context, id int) error {
				return errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := DeleteHandler{deleteFunc(tc.deleteFunc)}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type deleteFunc func(ctx context.Context, id int) error

func (d deleteFunc) Delete(ctx context.Context, id int) error {
	return d(ctx, id)
}
