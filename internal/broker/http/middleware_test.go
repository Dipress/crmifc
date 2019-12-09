package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/kit/auth"
)

func Test_authMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		header    map[string]string
		parseFunc func(ctx context.Context, tknStr string) (auth.Claims, error)
		code      int
	}{
		{
			name: "ok",
			header: map[string]string{
				"Authorization": "Bearer token",
			},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, nil
			},
			code: http.StatusOK,
		},
		{
			name:   "missing",
			header: map[string]string{},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, nil
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "wrong format",
			header: map[string]string{
				"Authorization": "token",
			},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, nil
			},
			code: http.StatusUnauthorized,
		},
		{
			name: "wrong token",
			header: map[string]string{
				"Authorization": "Bearer wrong",
			},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, errors.New("mock error")
			},
			code: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			next := handler.Func(func(w http.ResponseWriter, r *http.Request) error {
				return nil
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://exapmle.com", nil)

			for k, v := range tc.header {
				r.Header.Set(k, v)
			}

			authMiddleware(parseFunc(tc.parseFunc))(next).Handle(w, r)

			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected: %d", w.Code, tc.code)
			}
		})
	}
}

type parseFunc func(ctx context.Context, tknStr string) (auth.Claims, error)

func (p parseFunc) ParseClaims(ctx context.Context, tknStr string) (auth.Claims, error) {
	return p(ctx, tknStr)
}

func Test_contentTypeMiddleware(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://exapmle.com", nil)
	rec := httptest.NewRecorder()

	next := handler.Func(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	contentTypeMiddleware(next).Handle(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected to set application/json Content-Type header: %s", ct)
	}
}
