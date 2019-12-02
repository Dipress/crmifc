package role

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=handler.go -package=role -destination=handler.mock.go Service

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Service contains all services.
type Service interface {
	Create(ctx context.Context, f *role.Form) (*role.Role, error)
	Find(ctx context.Context, id int) (*role.Role, error)
	Update(ctx context.Context, id int, f *role.Form) (*role.Role, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) (*role.Roles, error)
}

// CreateHandler for create requests.
type CreateHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f role.Form

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	role, err := h.Create(r.Context(), &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "create role")
		}
	}

	data, err = role.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// FindHandler for find requests.
type FindHandler struct {
	Service
}

// Handle implements Handler interface.
func (rol *FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	rl, err := rol.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case role.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find role")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find role")
		}
	}

	data, err := rl.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// UpdateHandler for update requests.
type UpdateHandler struct {
	Service
}

// Handle implements Handler interface.
func (rol *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f role.Form
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	rl, err := rol.Update(r.Context(), id, &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "update role")
		}
	}

	data, err = rl.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// DeleteHandler for delete requests.
type DeleteHandler struct {
	Service
}

// Handle implements Handler interface.
func (rol *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := rol.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete role")
	}

	return nil
}

// ListHandler for list requests.
type ListHandler struct {
	Service
}

// Handle implements Handler interface.
func (rol *ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	roles, err := rol.List(r.Context())
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "list of roles")
	}

	data, err := roles.MarshalJSON()
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
