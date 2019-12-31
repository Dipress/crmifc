package article

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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
	List(ctx context.Context) (*article.Articles, error)
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
		return fmt.Errorf("read body: %w", response.BadRequestResponse(w))
	}

	if err := f.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal json: %w", response.BadRequestResponse(w))
	}

	article, err := h.Create(r.Context(), &f)
	if err != nil {
		var vErr validation.Errors
		switch {
		case errors.As(err, vErr):
			return fmt.Errorf("validation response: %w", response.UnprocessabeEntityResponse(w, vErr))
		default:
			return fmt.Errorf("create article: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err = article.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	a, err := f.Find(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, article.ErrNotFound):
			return fmt.Errorf("find: %w", response.NotFoundResponse(w))
		default:
			return fmt.Errorf("find article: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err := a.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", response.BadRequestResponse(w))
	}

	var f article.Form
	if err := f.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("unmarshal json: %w", response.BadRequestResponse(w))
	}

	art, err := h.Update(r.Context(), id, &f)
	if err != nil {
		var vErr validation.Errors
		switch {
		case errors.As(err, vErr):
			return fmt.Errorf("validation response: %w", response.UnprocessabeEntityResponse(w, vErr))
		default:
			return fmt.Errorf("update article: %w", response.InternalServerErrorResponse(w))
		}
	}

	data, err = art.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json: %w", response.InternalServerErrorResponse(w))
	}

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write response: %w", response.InternalServerErrorResponse(w))
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
		return fmt.Errorf("convert id query param to int: %v: %w", err, response.BadRequestResponse(w))
	}

	if err := h.Delete(r.Context(), id); err != nil {
		return fmt.Errorf("delete article: %w", response.InternalServerErrorResponse(w))
	}

	return nil
}

// ListHandler for article list request.
type ListHandler struct {
	Service
}

// Handle implements Handler interface.
func (h ListHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.List(r.Context())
	if err != nil {
		return fmt.Errorf("list of articles: %w", response.InternalServerErrorResponse(w))
	}

	data, err := articles.MarshalJSON()
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
