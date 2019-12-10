package handler

import (
	"net/http"
)

// Handler allows to handle requests.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

// Func type is an adapter to allow the use of
// ordinary functions as HTTP handlers
type Func func(http.ResponseWriter, *http.Request) error

// Handle implements Handler interface
func (f Func) Handle(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// Middleware represents middleware type.
type Middleware func(Handler) Handler

// Chain acts as a list of handler.Handler constructors.
type Chain struct {
	middlewares []Middleware
}

// NewChain creates new chain.
func NewChain(middlewares ...Middleware) Chain {
	return Chain{append(([]Middleware)(nil), middlewares...)}
}

// Then chains the middleware and returns the final handler.Handler.
func (c Chain) Then(h Handler) Handler {
	if len(c.middlewares) == 0 {
		return h
	}

	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}

	return h
}

// Append extends a chain, adding the specified middlewares.
func (c Chain) Append(middlewares ...Middleware) Chain {
	newMiddlewares := make([]Middleware, 0, len(c.middlewares)+len(middlewares))
	newMiddlewares = append(newMiddlewares, c.middlewares...)
	newMiddlewares = append(newMiddlewares, middlewares...)

	return Chain{newMiddlewares}
}
