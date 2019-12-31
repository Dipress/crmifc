package http

import (
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
	"github.com/dipress/crmifc/internal/kit/logger"
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
func NewServer(addr string, logger *logger.Logger, services *Services, authenticator *authEng.Authenticator) *http.Server {
	mux := mux.NewRouter().StrictSlash(true)

	// Auth handler.
	authenticateHandler := authHandlers.AuthenticaterHandler{
		Authenticater: services.Auth,
	}

	base := handler.NewChain(contentTypeMiddleware)

	// Auth route.
	mux.Handle("/signin", finalizeMiddleware(logger, base)(&authenticateHandler)).Methods(http.MethodPost)

	authorized := base.Append(authMiddleware(authenticator))

	articles := mux.PathPrefix("/articles").Subrouter()
	articleHandlers.Prepare(articles, services.Article, finalizeMiddleware(logger, authorized))

	categories := mux.PathPrefix("/categories").Subrouter()
	categoryHandlers.Prepare(categories, services.Category, finalizeMiddleware(logger, authorized))

	admin := authorized.Append(adminMiddleware(abillity.UserAbillity{}))

	roles := mux.PathPrefix("/roles").Subrouter()
	roleHandlers.Prepare(roles, services.Role, finalizeMiddleware(logger, admin))

	users := mux.PathPrefix("/users").Subrouter()
	userHandlers.Prepare(users, services.User, finalizeMiddleware(logger, admin))

	// TODO: Delete after deploy to real server.
	// allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	// allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"})
	// allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	s := http.Server{
		Addr:         addr,
		Handler:      mux, //handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(mux),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	return &s
}
