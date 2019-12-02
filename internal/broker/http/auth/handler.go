package auth

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/pkg/errors"
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
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	if err := a.Authenticater.Authenticate(r.Context(), f.Email, f.Password, &t); err != nil {
		switch err := errors.Cause(err); err {
		case auth.ErrEmailNotFound, auth.ErrWrongPassword:
			return errors.Wrap(response.UnauthorizedResponse(w), "find user")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "authenticate user")
		}
	}

	data, err = t.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}
