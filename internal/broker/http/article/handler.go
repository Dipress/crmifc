package article

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/broker/http/response"
	"github.com/dipress/crmifc/internal/validation"
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
			return errors.Wrap(response.InternalServerErrorResponse(w), "create role")
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
