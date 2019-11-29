package user

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/user"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
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

// Creater abstraction for create service.
type Creater interface {
	Create(ctx context.Context, f *user.Form, u *user.User) error
}

// Finder abstraction for find service.
type Finder interface {
	Find(ctx context.Context, id int) (*user.User, error)
}

// Updater abstraction for update service.
type Updater interface {
	Update(ctx context.Context, id int, f *user.Form) error
}

// Deleter abstraction for delete service.
type Deleter interface {
	Delete(ctx context.Context, id int) error
}

// Lister absctraction for list service.
type Lister interface {
	List(ctx context.Context) (*user.Users, error)
}

// AuthHandler for authenticate request.
type AuthHandler struct {
	Authenticater
}

// Handle implements Handler interface.
func (a AuthHandler) Handle(w http.ResponseWriter, r *http.Request) error {
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
			return errors.Wrap(response.InternalServerErrorResponse(w), "authenticate")
		}
	}

	data, err = t.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// CreateHandler for  user create requests.
type CreateHandler struct {
	Creater
}

// Handle implements Handler interface.
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f user.Form

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	var u user.User
	if err := h.Creater.Create(r.Context(), &f, &u); err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "registrate")
		}
	}

	data, err = u.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// FindHandler for  user create requests.
type FindHandler struct {
	Finder
}

// Handle implements Handler interface.
func (h FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	u, err := h.Finder.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find")
		}
	}

	data, err := u.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// UpdateHandler for  user update requests.
type UpdateHandler struct {
	Updater
}

// Handle implements Handler interface.
func (h UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	var f user.Form
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	if err := h.Updater.Update(r.Context(), id, &f); err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "registrate")
		}
	}

	return nil
}

// DeleteHandler for user delete requests.
type DeleteHandler struct {
	Deleter
}

// Handle implements Handler interface.
func (u *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := u.Deleter.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete role")
	}

	return nil
}

// ListHandler for user list request.
type ListHandler struct {
	Lister
}

// Handle implements Handler interface.
func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	users, err := h.Lister.List(r.Context())
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "list of users")
	}

	data, err := users.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}
