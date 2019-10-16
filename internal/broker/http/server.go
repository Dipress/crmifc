package http

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/gorilla/mux"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/broker/http/category"
	"github.com/dipress/crmifc/internal/broker/http/role"
	"github.com/dipress/crmifc/internal/broker/http/user"
	categoryCreate "github.com/dipress/crmifc/internal/category/create"
	authEng "github.com/dipress/crmifc/internal/kit/auth"
	roleCreate "github.com/dipress/crmifc/internal/role/create"
	roleDelete "github.com/dipress/crmifc/internal/role/delete"
	roleFind "github.com/dipress/crmifc/internal/role/find"
	roleList "github.com/dipress/crmifc/internal/role/list"
	roleUpdate "github.com/dipress/crmifc/internal/role/update"
	userCreate "github.com/dipress/crmifc/internal/user/create"
	userDelete "github.com/dipress/crmifc/internal/user/delete"
	userFind "github.com/dipress/crmifc/internal/user/find"
	userList "github.com/dipress/crmifc/internal/user/list"
	userUpdate "github.com/dipress/crmifc/internal/user/update"
)

const (
	timeout = 15 * time.Second
)

// NewServer prepare http server to work.
func NewServer(addr string, db *sql.DB, authenticator *authEng.Authenticator) *http.Server {
	mux := mux.NewRouter()

	repo := postgres.NewRepository(db)
	authenticateService := auth.NewService(repo, authenticator, time.Hour*24)
	rolesCreateService := roleCreate.NewService(repo, &validation.Role{})
	roleFindService := roleFind.NewService(repo)
	roleUpdateService := roleUpdate.NewService(repo, &validation.Role{})
	roleDeleteService := roleDelete.NewService(repo)
	roleListService := roleList.NewService(repo)

	userCreateService := userCreate.NewService(repo, &validation.User{})
	userFindService := userFind.NewService(repo)
	userUpdateService := userUpdate.NewService(repo, &validation.User{})
	userDeleteService := userDelete.NewService(repo)
	userListService := userList.NewService(repo)

	categoryCreateService := categoryCreate.NewService(repo, &validation.Category{})

	// Auth handler.
	authenticateHandler := user.AuthHandler{
		Authenticater: authenticateService,
	}

	// Auth route.
	mux.HandleFunc("/signin", user.HTTPHandler{
		Handler: &authenticateHandler,
	}.ServeHTTP).Methods(http.MethodPost)

	// User handlers.
	userCreateHandler := user.CreateHandler{
		Creater: userCreateService,
	}

	userFindHandler := user.FindHandler{
		Finder: userFindService,
	}

	userUpdateHandler := user.UpdateHandler{
		Updater: userUpdateService,
	}

	userDeleteHandler := user.DeleteHandler{
		Deleter: userDeleteService,
	}

	userListHandler := user.ListHandler{
		Lister: userListService,
	}

	// Role handlers.
	roleCreateHandler := role.CreateHandler{
		Creater: rolesCreateService,
	}

	roleFindHandler := role.FindHandler{
		Finder: roleFindService,
	}

	roleUpdateHandler := role.UpdateHandler{
		Updater: roleUpdateService,
	}

	roleDeleteHandler := role.DeleteHandler{
		Deleter: roleDeleteService,
	}

	roleListHandler := role.ListHanlder{
		Lister: roleListService,
	}

	// Category handlers.
	categoryCreateHandler := category.CreateHandler{
		Creater: categoryCreateService,
	}

	// Category routes.
	mux.HandleFunc("/categories", AuthMiddleware(role.HTTPHandler{
		Handler: &categoryCreateHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodPost)

	// User routes.
	mux.HandleFunc("/users", AuthMiddleware(user.HTTPHandler{
		Handler: &userCreateHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodPost)

	mux.HandleFunc("/users/{id}", AuthMiddleware(user.HTTPHandler{
		Handler: &userFindHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodGet)

	mux.HandleFunc("/users/{id}", AuthMiddleware(user.HTTPHandler{
		Handler: &userUpdateHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodPut)

	mux.HandleFunc("/users/{id}", AuthMiddleware(user.HTTPHandler{
		Handler: &userDeleteHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodDelete)

	mux.HandleFunc("/users", AuthMiddleware(user.HTTPHandler{
		Handler: &userListHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodGet)

	// Role routes.
	mux.HandleFunc("/roles", AuthMiddleware(role.HTTPHandler{
		Handler: &roleCreateHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodPost)

	mux.HandleFunc("/roles/{id}", AuthMiddleware(role.HTTPHandler{
		Handler: &roleFindHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodGet)

	mux.HandleFunc("/roles/{id}", AuthMiddleware(role.HTTPHandler{
		Handler: &roleUpdateHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodPut)

	mux.HandleFunc("/roles/{id}", AuthMiddleware(role.HTTPHandler{
		Handler: &roleDeleteHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodDelete)

	mux.HandleFunc("/roles", AuthMiddleware(role.HTTPHandler{
		Handler: &roleListHandler,
	}, authenticator).ServeHTTP).Methods(http.MethodGet)

	s := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	return &s
}
