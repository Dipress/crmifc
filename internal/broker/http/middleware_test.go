package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dipress/crmifc/internal/kit/auth"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		header    map[string]string
		parseFunc func(ctx context.Context, tknStr string) (auth.Claims, error)
		callNext  bool
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
			callNext: true,
			code:     http.StatusOK,
		},
		{
			name:   "missing",
			header: map[string]string{},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, nil
			},
			callNext: false,
			code:     http.StatusUnauthorized,
		},
		{
			name: "wrong format",
			header: map[string]string{
				"Authorization": "token",
			},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, nil
			},
			callNext: false,
			code:     http.StatusUnauthorized,
		},
		{
			name: "wrong token",
			header: map[string]string{
				"Authorization": "Bearer wrong",
			},
			parseFunc: func(ctx context.Context, tknStr string) (auth.Claims, error) {
				return auth.Claims{}, errors.New("mock error")
			},
			callNext: false,
			code:     http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			nextCalls := make(chan struct{})
			b := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				go func() {
					nextCalls <- struct{}{}
				}()
			})
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://exapmle.com", nil)
			for k, v := range tc.header {
				r.Header.Set(k, v)
			}
			h := AuthMiddleware(b, parseFunc(tc.parseFunc))
			h.ServeHTTP(w, r)

			if tc.callNext {
				select {
				case <-nextCalls:
				case <-time.After(time.Second):
					t.Error("should write to next channel")
				}
				return
			}

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
