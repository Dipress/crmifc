package user

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/user"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=handler.go -package=user -destination=handler.mock.go Service

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Service contains all services.
type Service interface {
	Create(ctx context.Context, f *user.Form, u *user.User) error
	Find(ctx context.Context, id int) (*user.User, error)
	Update(ctx context.Context, id int, f *user.Form) (*user.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) (*user.Users, error)
}

// CreateHandler for  user create requests.
type CreateHandler struct {
	Service
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
	if err := h.Create(r.Context(), &f, &u); err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "create user")
		}
	}

	data, err = u.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// FindHandler for  user create requests.
type FindHandler struct {
	Service
}

// Handle implements Handler interface.
func (h FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	u, err := h.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case user.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find user")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find user")
		}
	}

	data, err := u.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// UpdateHandler for  user update requests.
type UpdateHandler struct {
	Service
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

	u, err := h.Update(r.Context(), id, &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "update user")
		}
	}

	data, err = u.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// DeleteHandler for user delete requests.
type DeleteHandler struct {
	Service
}

// Handle implements Handler interface.
func (u *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := u.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete user")
	}

	return nil
}

// ListHandler for user list request.
type ListHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	users, err := h.List(r.Context())
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "list of users")
	}

	data, err := users.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// Prepare prepares routes to use.
func Prepare(subrouter *mux.Router, service Service, middleware func(handler.Handler) http.Handler) {
	create := CreateHandler{service}
	find := FindHandler{service}
	update := UpdateHandler{service}
	delete := DeleteHandler{service}
	list := ListHandler{service}

	subrouter.Handle("", middleware(&create)).Methods(http.MethodPost)
	subrouter.Handle("/{id}", middleware(&find)).Methods(http.MethodGet)
	subrouter.Handle("/{id}", middleware(&update)).Methods(http.MethodPut)
	subrouter.Handle("/{id}", middleware(&delete)).Methods(http.MethodDelete)
	subrouter.Handle("", middleware(&list)).Methods(http.MethodGet)
}
