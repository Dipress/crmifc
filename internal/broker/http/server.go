package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dipress/crmifc/internal/abillity"
	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/auth"
	articleHandlers "github.com/dipress/crmifc/internal/broker/http/article"
	authHandlers "github.com/dipress/crmifc/internal/broker/http/auth"
	categoryHandlers "github.com/dipress/crmifc/internal/broker/http/category"
	"github.com/dipress/crmifc/internal/broker/http/handler"
	roleHandlers "github.com/dipress/crmifc/internal/broker/http/role"
	userHandlers "github.com/dipress/crmifc/internal/broker/http/user"
	"github.com/dipress/crmifc/internal/category"
	authEng "github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/user"
)

const (
	timeout = 15 * time.Second
)

// Services contains all the services.
type Services struct {
	Auth     *auth.Service
	Article  *article.Service
	Category *category.Service
	Role     *role.Service
	User     *user.Service
}

// NewServer prepare http server to work.
func NewServer(addr string, services *Services, authenticator *authEng.Authenticator) *http.Server {
	mux := mux.NewRouter().StrictSlash(true)

	// Auth handler.
	authenticateHandler := authHandlers.AuthenticaterHandler{
		Authenticater: services.Auth,
	}

	base := handler.NewChain(contentTypeMiddleware)

	// Auth route.
	mux.Handle("/signin", finalizeMiddleware(base)(&authenticateHandler)).Methods(http.MethodPost)

	authorized := base.Append(authMiddleware(authenticator))

	articles := mux.PathPrefix("/articles").Subrouter()
	articleHandlers.Prepare(articles, services.Article, finalizeMiddleware(authorized))

	categories := mux.PathPrefix("/categories").Subrouter()
	categoryHandlers.Prepare(categories, services.Category, finalizeMiddleware(authorized))

	admin := authorized.Append(adminMiddleware(abillity.UserAbillity{}))

	roles := mux.PathPrefix("/roles").Subrouter()
	roleHandlers.Prepare(roles, services.Role, finalizeMiddleware(admin))

	users := mux.PathPrefix("/users").Subrouter()
	userHandlers.Prepare(users, services.User, finalizeMiddleware(admin))

	s := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	return &s
}

func finalizeMiddleware(middleware handler.Chain) func(handler.Handler) http.Handler {
	f := func(handler handler.Handler) http.Handler {
		wrapped := middleware.Then(handler)
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := wrapped.Handle(w, r); err != nil {
				log.Printf("serve http: %+v\n", err)
			}
		})

		return h
	}

	return f
}
