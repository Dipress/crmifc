package role

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
)

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		createFunc func(ctx context.Context, f *role.Form) (*role.Role, error)
		code       int
	}{
		{
			name: "ok",
			createFunc: func(ctx context.Context, f *role.Form) (*role.Role, error) {
				return &role.Role{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			createFunc: func(ctx context.Context, f *role.Form) (*role.Role, error) {
				return &role.Role{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			createFunc: func(ctx context.Context, f *role.Form) (*role.Role, error) {
				return &role.Role{}, errors.New("mock error")
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

type createFunc func(ctx context.Context, f *role.Form) (*role.Role, error)

func (c createFunc) Create(ctx context.Context, f *role.Form) (*role.Role, error) {
	return c(ctx, f)
}

func TestFindHandler(t *testing.T) {
	tests := []struct {
		name     string
		findFunc func(ctx context.Context, id int) (*role.Role, error)
		code     int
	}{
		{
			name: "ok",
			findFunc: func(ctx context.Context, id int) (*role.Role, error) {
				return &role.Role{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "internal error",
			findFunc: func(ctx context.Context, id int) (*role.Role, error) {
				return &role.Role{}, errors.New("mock error")
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

type findFunc func(ctx context.Context, id int) (*role.Role, error)

func (f findFunc) Find(ctx context.Context, id int) (*role.Role, error) {
	return f(ctx, id)
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		updateFunc func(ctx context.Context, id int, f *role.Form) (*role.Role, error)
		code       int
	}{
		{
			name: "ok",
			updateFunc: func(ctx context.Context, id int, f *role.Form) (*role.Role, error) {
				return &role.Role{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			updateFunc: func(ctx context.Context, id int, f *role.Form) (*role.Role, error) {
				return &role.Role{}, make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			updateFunc: func(ctx context.Context, id int, f *role.Form) (*role.Role, error) {
				return &role.Role{}, errors.New("mock error")
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

type updateFunc func(ctx context.Context, id int, f *role.Form) (*role.Role, error)

func (u updateFunc) Update(ctx context.Context, id int, f *role.Form) (*role.Role, error) {
	return u(ctx, id, f)
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

func TestListHandler(t *testing.T) {
	tests := []struct {
		name     string
		listFunc func(ctx context.Context) (*role.Roles, error)
		code     int
	}{
		{
			name: "ok",
			listFunc: func(ctx context.Context) (*role.Roles, error) {
				return &role.Roles{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "repository error",
			listFunc: func(ctx context.Context) (*role.Roles, error) {
				return &role.Roles{}, errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := ListHandler{listFunc(tc.listFunc)}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type listFunc func(ctx context.Context) (*role.Roles, error)

func (l listFunc) List(ctx context.Context) (*role.Roles, error) {
	return l(ctx)
}
