package article

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/broker/http/response"
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
	Create(ctx context.Context, f *article.Form) (*article.Article, error)
}

// Finder abstraction for find service.
type Finder interface {
	Find(ctx context.Context, id int) (*article.Article, error)
}

// Updater abstraction for update service.
type Updater interface {
	Update(ctx context.Context, id int, f *article.Form) (*article.Article, error)
}

// Deleter abstraction for delete service.
type Deleter interface {
	Delete(ctx context.Context, id int) error
}

// CreateHandler for create requests.
type CreateHandler struct {
	Creater
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

	article, err := h.Creater.Create(r.Context(), &f)
	if err != nil {
		fmt.Println(err)
		switch v := errors.Cause(err).(type) {
		case validation.Errors:
			return errors.Wrap(response.UnprocessabeEntityResponse(w, v), "validation response")
		default:
			return errors.Wrap(response.InternalServerErrorResponse(w), "create article")
		}
	}

	data, err = article.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// FindHandler for article create requests.
type FindHandler struct {
	Finder
}

// Handle implements Handler interface.
func (f FindHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	a, err := f.Finder.Find(r.Context(), id)
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
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// UpdateHandler for article update requests.
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "read body")
	}

	var f article.Form
	if err := f.UnmarshalJSON(data); err != nil {
		return errors.Wrap(response.BadRequestResponse(w), "unmarshal json")
	}

	art, err := h.Updater.Update(r.Context(), id, &f)
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
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// DeleteHandler for article update requests.
type DeleteHandler struct {
	Deleter
}

// Handle implements Handler interface.
func (h DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.Wrapf(response.BadRequestResponse(w), "convert id query param to int: %v", err)
	}

	if err := h.Deleter.Delete(r.Context(), id); err != nil {
		return errors.Wrap(response.InternalServerErrorResponse(w), "delete article")
	}

	return nil
}
