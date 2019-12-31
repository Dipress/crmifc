package role

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
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
		return fmt.Errorf("read body: %w", response.BadRequestResponse(w))
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal json: %w", response.BadRequestResponse(w))
	}

	role, err := h.Create(r.Context(), &f)
	if err != nil {
		var vErr validation.Errors
		switch {
		case errors.As(err, vErr):
			return fmt.Errorf("validation response: %w", response.UnprocessabeEntityResponse(w, vErr))
		default:
			return fmt.Errorf("create role: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err = role.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	rl, err := rol.Find(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, role.ErrNotFound):
			return fmt.Errorf("find role: %w", response.NotFoundResponse(w))
		default:
			return fmt.Errorf("find role: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err := rl.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", response.BadRequestResponse(w))
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal json: %w", response.BadRequestResponse(w))
	}

	rl, err := rol.Update(r.Context(), id, &f)
	if err != nil {
		var vErr validation.Errors
		switch {
		case errors.As(err, vErr):
			return fmt.Errorf("validation response: %w", response.UnprocessabeEntityResponse(w, vErr))
		default:
			return fmt.Errorf("update role: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err = rl.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	if err := rol.Delete(r.Context(), id); err != nil {
		return fmt.Errorf("delete role: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("list of roles: %w", response.InternalServerErrorResponse(w))
	}

	data, err := roles.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
