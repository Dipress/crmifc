package http

import (
	"fmt"
	"net/http"

	"github.com/dipress/crmifc/internal/broker/http/handler"
)

// Logger allows to log info and errors.
type Logger interface {
	Info(info string, extra map[string]interface{})
	Warn(info string, extra map[string]interface{})
	Error(err error, extra map[string]interface{})
	Fatal(err error, extra map[string]interface{})
}

func finalizeMiddleware(logger Logger, middleware handler.Chain) func(handler.Handler) http.Handler {
	f := func(handler handler.Handler) http.Handler {
		wrapped := middleware.Then(handler)
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := wrapped.Handle(w, r); err != nil {
				logger.Error(fmt.Errorf("serve http error: %w", err), map[string]interface{}{
					"path":   r.URL.Path,
					"method": r.Method,
				})
			}
		})

		return h
	}

	return f
}
