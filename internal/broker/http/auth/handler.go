package auth

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/broker/http/response"
)

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Authenticater abstraction for authenticate service.
type Authenticater interface {
	Authenticate(ctx context.Context, email, password string, t *auth.Token) error
}

// AuthenticaterHandler for authenticate request.
type AuthenticaterHandler struct {
	Authenticater
}

// Handle implements Handler interface.
func (a AuthenticaterHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f auth.Form
	var t auth.Token

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", response.BadRequestResponse(w))
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal json: %w", response.BadRequestResponse(w))
	}

	if err := a.Authenticater.Authenticate(r.Context(), f.Email, f.Password, &t); err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailNotFound), errors.Is(err, auth.ErrWrongPassword):
			return fmt.Errorf("find user: %w", response.UnauthorizedResponse(w))
		default:
			return fmt.Errorf("authenticate user: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err = t.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
	}

	return nil
}
