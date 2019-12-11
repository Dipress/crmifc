package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dipress/crmifc/internal/broker/http/response"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/kit/auth"
)

// Authenticator is used to authenticate clients.
// It recreates the claims by parsing the token.
type Authenticator interface {
	ParseClaims(ctx context.Context, tknStr string) (auth.Claims, error)
}

// authMiddleware represents middleware with authentication.
func authMiddleware(a Authenticator) handler.Middleware {
	m := func(next handler.Handler) handler.Handler {
		h := handler.Func(func(w http.ResponseWriter, r *http.Request) error {
			authHdr := r.Header.Get("Authorization")
			if authHdr == "" {
				return response.UnauthorizedResponse(w)
			}

			tknStr, err := parseAuthHeader(authHdr)
			if err != nil {
				return response.UnauthorizedResponse(w)
			}

			c := r.Context()
			cl, err := a.ParseClaims(c, tknStr)
			if err != nil {
				return response.UnauthorizedResponse(w)
			}

			ctx := auth.ToContext(c, &cl)
			r = r.WithContext(ctx)

			return next.Handle(w, r)
		})

		return h
	}

	return m
}

// contentTypeMiddleware sets content type header.
func contentTypeMiddleware(next handler.Handler) handler.Handler {
	h := handler.Func(func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		return next.Handle(w, r)
	})

	return h
}

// parseAuthHeader parses an authorization header. Expected header is of
// the format `Bearer <token>`.
func parseAuthHeader(bearerStr string) (string, error) {
	split := strings.Split(bearerStr, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("expected Authorization header format: Bearer <token>")
	}

	return split[1], nil
}
