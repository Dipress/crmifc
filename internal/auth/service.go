package auth

import (
	"context"
	"time"

	"github.com/dipress/blog/kit/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/dipress/crmifc/internal/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// easyjson service.go

var (
	// ErrEmailNotFound returns when given email is not
	// found in database.
	ErrEmailNotFound = errors.New("email not found")
	// ErrWrongPassword returns when given password
	// isn't equal to to its hash in the database.
	ErrWrongPassword = errors.New("wrong password")
)

// Repository allows working with a database.
type Repository interface {
	FindByEmail(ctx context.Context, email string, user *user.User) error
}

// TokenGenerator generates token for authenticated user.
type TokenGenerator interface {
	GenerateToken(ctx context.Context, claims jwt.Claims) (string, error)
}

// Service holds required data for user
// authentication.
type Service struct {
	Repository
	TokenGenerator
	ExpireAfter time.Duration
}

// Form is a user auth form.
//easyjson:json
type Form struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token holds token data.
//easyjson:json
type Token struct {
	Token string `json:"token"`
}

// NewService factory takes in required arguments
// and returns a pointer to the Service instance.
func NewService(r Repository, t TokenGenerator, exp time.Duration) *Service {
	s := Service{
		Repository:     r,
		TokenGenerator: t,
		ExpireAfter:    exp,
	}

	return &s
}

// Authenticate allows authenticating user by given email and password
// and set t Token value as generated token.
func (s *Service) Authenticate(ctx context.Context, email, password string, t *Token) error {
	var user user.User
	if err := s.Repository.FindByEmail(ctx, email, &user); err != nil {
		return errors.Wrap(err, "find user by email")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return ErrWrongPassword
	}

	// If we are this far the request is valid.
	// Now we need to create the token for the user.
	claims := auth.NewClaims(user.Email, time.Now(), s.ExpireAfter)

	tknStr, err := s.GenerateToken(ctx, claims)
	if err != nil {
		return errors.Wrap(err, "generate token")
	}
	t.Token = tknStr

	return nil
}
