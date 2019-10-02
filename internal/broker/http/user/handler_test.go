package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/user"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
)

func TestAuthHandler(t *testing.T) {
	tests := []struct {
		name     string
		authFunc func(ctx context.Context, email, password string, t *auth.Token) error
		code     int
	}{
		{
			name: "ok",
			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "email error",
			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return auth.ErrEmailNotFound
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "password error",
			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return auth.ErrWrongPassword
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "internal error",
			authFunc: func(ctx context.Context, email, password string, t *auth.Token) error {
				return errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := AuthHandler{authFunc(tc.authFunc)}
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

type authFunc func(ctx context.Context, email, password string, t *auth.Token) error

func (a authFunc) Authenticate(ctx context.Context, email, password string, t *auth.Token) error {
	return a(ctx, email, password, t)
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name       string
		createFunc func(ctx context.Context, f *user.Form, user *user.User) error
		code       int
	}{
		{
			name: "ok",
			createFunc: func(ctx context.Context, f *user.Form, user *user.User) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			createFunc: func(ctx context.Context, f *user.Form, user *user.User) error {
				return make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			createFunc: func(ctx context.Context, f *user.Form, user *user.User) error {
				return errors.New("mock error")
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

type createFunc func(ctx context.Context, f *user.Form, user *user.User) error

func (c createFunc) Create(ctx context.Context, f *user.Form, user *user.User) error {
	return c(ctx, f, user)
}

func TestFindHandler(t *testing.T) {
	tests := []struct {
		name     string
		findFunc func(ctx context.Context, id int) (*user.User, error)
		code     int
	}{
		{
			name: "ok",
			findFunc: func(ctx context.Context, id int) (*user.User, error) {
				return &user.User{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "internal error",
			findFunc: func(ctx context.Context, id int) (*user.User, error) {
				return &user.User{}, errors.New("mock error")
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

type findFunc func(ctx context.Context, id int) (*user.User, error)

func (f findFunc) Find(ctx context.Context, id int) (*user.User, error) {
	return f(ctx, id)
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name       string
		updateFunc func(ctx context.Context, id int, f *user.Form) error
		code       int
	}{
		{
			name: "ok",
			updateFunc: func(ctx context.Context, id int, f *user.Form) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			updateFunc: func(ctx context.Context, id int, f *user.Form) error {
				return make(validation.Errors)
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			updateFunc: func(ctx context.Context, id int, f *user.Form) error {
				return errors.New("mock error")
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

type updateFunc func(ctx context.Context, id int, f *user.Form) error

func (u updateFunc) Update(ctx context.Context, id int, f *user.Form) error {
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
			r := httptest.NewRequest("PUT", "http://example.com", strings.NewReader("{}"))
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
		listFunc func(ctx context.Context) (*user.Users, error)
		code     int
	}{
		{
			name: "ok",
			listFunc: func(ctx context.Context) (*user.Users, error) {
				return &user.Users{}, nil
			},
			code: http.StatusOK,
		},
		{
			name: "repository error",
			listFunc: func(ctx context.Context) (*user.Users, error) {
				return &user.Users{}, errors.New("mock error")
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
			r := httptest.NewRequest("GET", "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}

		})
	}
}

type listFunc func(ctx context.Context) (*user.Users, error)

func (l listFunc) List(ctx context.Context) (*user.Users, error) {
	return l(ctx)
}
