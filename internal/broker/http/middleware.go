package http

import (
	"context"
	"errors"
	"net/http"
	"os/user"
	"strings"

	"github.com/dipress/crmifc/internal/kit/auth"
)

// Authenticator is used to authenticate clients.
// It recreates the claims by parsing the token.
type Authenticator interface {
	ParseClaims(ctx context.Context, tknStr string) (auth.Claims, error)
}

// AuthMiddleware represents middleware with authentication.
func AuthMiddleware(next http.Handler, a Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		authHdr := r.Header.Get("Authorization")
		if authHdr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tknStr, err := parseAuthHeader(authHdr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cl, err := a.ParseClaims(c, tknStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := auth.ToContext(c, &cl)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
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

// Repository allows working with a database.
type Repository interface {
	FindByEmail(ctx context.Context, email string, usr *user.User) error
}

// PermissionMiddleware represents middleware with permissions.
func PermissionMiddleware(next http.Handler, r Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
