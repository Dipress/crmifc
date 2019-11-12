package article

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/broker/http/handler"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/validation"
)

// go:generate mockgen -source=handler.go -package=article -destination=handler.mock.go Service

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Service contains all services.
type Service interface {
	Create(ctx context.Context, f *article.Form) (*article.Article, error)
	Find(ctx context.Context, id int) (*article.Article, error)
	Update(ctx context.Context, id int, f *article.Form) (*article.Article, error)
	Delete(ctx context.Context, id int) error
}

// CreateHandler for create requests.
type CreateHandler struct {
	Service
}

// Handle implements Handler interface.
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var f article.Form

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	article, err := h.Create(r.Context(), &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "create article")
		}
	}

	data, err = article.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// FindHandler for article create requests.
type FindHandler struct {
	Service
}

// Handle implements Handler interface.
func (f FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	a, err := f.Find(r.Context(), id)
	if err != nil {
		switch errors.Cause(err) {
		case article.ErrNotFound:
			return errors.Wrap(response.NotFoundResponse(w), "find")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "find article")
		}
	}

	data, err := a.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// UpdateHandler for article update requests.
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	var f article.Form
	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	art, err := h.Update(r.Context(), id, &f)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "update article")
		}
	}

	data, err = art.MarshalJSON()
	if err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "write response")
	}

	return nil
}

// DeleteHandler for article update requests.
type DeleteHandler struct {
	Service
}

// Handle implements Handler interface.
func (h DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := h.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete article")
	}

	return nil
}

// Prepare prepares routes to use.
func Prepare(subrouter *mux.Router, service Service, middleware func(handler.Handler) http.Handler) {
	create := CreateHandler{service}
	find := FindHandler{service}
	update := UpdateHandler{service}
	delete := DeleteHandler{service}

	subrouter.Handle("", middleware(&create)).Methods(http.MethodPost)
	subrouter.Handle("/{id}", middleware(&find)).Methods(http.MethodGet)
	subrouter.Handle("/{id}", middleware(&update)).Methods(http.MethodPut)
	subrouter.Handle("/{id}", middleware(&delete)).Methods(http.MethodDelete)
}
