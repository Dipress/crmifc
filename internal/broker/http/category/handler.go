package category

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// go:generate mockgen -source=handler.go -package=category -destination=handler.mock.go Service

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Service contains all services.
type Service interface {
	Create(ctx context.Context, f *category.Form) (*category.Category, error)
	Find(ctx context.Context, id int) (*category.Category, error)
	Update(ctx context.Context, id int, f *category.Form) (*category.Category, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) (*category.Categories, error)
}

// CreateHandler for create requests.
type CreateHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f category.Form

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	category, err := h.Create(r.Context(), &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "create category")
		}
	}

	data, err = category.MarshalJSON()
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
func (h *FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	cat, err := h.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case category.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find category")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find category")
		}
	}

	data, err := cat.MarshalJSON()
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
func (h *UpdateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f category.Form
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

	cat, err := h.Update(r.Context(), id, &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "update category")
		}
	}

	data, err = cat.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}
	return nil
}

// DeleteHandler for delete request.
type DeleteHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := h.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete category")
	}

	return nil
}

// ListHandler for list requests.
type ListHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.List(r.Context())
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "list of categories")
	}

	data, err := categories.MarshalJSON()
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
