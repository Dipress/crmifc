package role

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// HTTPHandler allows to implement ServeHTTP for Handler.
type HTTPHandler struct {
	Handler
}

// ServeHTTP implements http.Handler.
func (h HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.Handle(w, r); err != nil {
		log.Printf("serve http: %+v\n", err)
	}
}

// Creater abstraction for create service.
type Creater interface {
	Create(ctx context.Context, f *role.Form) (*role.Role, error)
}

// Finder abstraction for find service.
type Finder interface {
	Find(ctx context.Context, id int) (*role.Role, error)
}

// Updater abstraction for update service.
type Updater interface {
	Update(ctx context.Context, id int, f *role.Form) (*role.Role, error)
}

// Deleter abstraction for delete service.
type Deleter interface {
	Delete(ctx context.Context, id int) error
}

// Lister abstraction for list service.
type Lister interface {
	List(ctx context.Context) (*role.Roles, error)
}

// CreateHandler for create requests.
type CreateHandler struct {
	Creater
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

	role, err := h.Creater.Create(r.Context(), &f)
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
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// FindHandler for find requests.
type FindHandler struct {
	Finder
}

// Handle implements Handler interface.
func (rol *FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	rl, err := rol.Finder.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case role.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find")
		}
	}

	data, err := rl.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// UpdateHandler for update requests.
type UpdateHandler struct {
	Updater
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

	rl, err := rol.Updater.Update(r.Context(), id, &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "update")
		}
	}

	data, err = rl.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// DeleteHandler for delete requests.
type DeleteHandler struct {
	Deleter
}

// Handle implements Handler interface.
func (rol *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := rol.Deleter.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete role")
	}

	return nil
}

// ListHanlder for list requests.
type ListHandler struct {
	Lister
}

// Handle implements Handler interface.
func (rol *ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	roles, err := rol.Lister.List(r.Context())
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "list of roles")
	}

	data, err := roles.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}
	return nil
}
