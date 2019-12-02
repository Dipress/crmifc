package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dipress/crmifc/internal/auth"
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
			h := AuthenticaterHandler{authFunc(tc.authFunc)}
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
