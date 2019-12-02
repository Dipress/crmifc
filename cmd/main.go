package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dipress/crmifc/internal/article"
	authSrv "github.com/dipress/crmifc/internal/auth"
	httpBroker "github.com/dipress/crmifc/internal/broker/http"
	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/storage/postgres/schema"
	"github.com/dipress/crmifc/internal/user"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/mattes/migrate"
	"github.com/pkg/errors"
)

const (
	alg = "RS256"
)

func main() {
	var (
		addr           = flag.String("addr", ":8080", "address of http server")
		dsn            = flag.String("dsn", "", "postgres database DSN")
		privateKeyFile = flag.String("key", "./internal/kit/keys/jwtRS256.key", "private key file path")
		keyID          = flag.String("id", "123456", "private key id")
	)
	flag.Parse()

	// Setup database connection.
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v\n", err)
	}
	defer db.Close()

	// Migrate schema.
	if err := schema.Migrate(db); err != nil {
		if errors.Cause(err) != migrate.ErrNoChange {
			log.Fatalf("failed to migrate schema: %v", err)
		}
	}

	// Seed data.
	if err := schema.Seed(db); err != nil {
		log.Fatalf("failed to seeds data: %v", err)
	}

	// Reppositories.
	userRepo := postgres.NewUserRepository(db)

	// Authentication setup.
	keyContents, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatalf("reading auth private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatalf("parsing auth private key: %v", err)
	}
	publicKeyLookup := auth.NewSingleKeyFunc(*keyID, key.Public().(*rsa.PublicKey))
	authenticator, err := auth.NewAuthenticator(key, *keyID, alg, publicKeyLookup, userRepo)
	if err != nil {
		log.Fatalf("constructing authenticator: %v", err)
	}

	// Services
	services := setupServices(db, authenticator)

	// Setup server.
	srv := setupServer(*addr, services, authenticator)
	if err := srv.ListenAndServe(); err != nil {
		errors.Wrap(err, "filed to serve http")
	}
}

func setupServer(addr string, services *httpBroker.Services, authenticator *auth.Authenticator) *http.Server {
	return httpBroker.NewServer(addr, services, authenticator)
}

func setupServices(db *sql.DB, authenticator *auth.Authenticator) *httpBroker.Services {
	// Repositorires
	userRepo := postgres.NewUserRepository(db)
	articleRepo := postgres.NewArticleRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	roleRepo := postgres.NewRoleRepository(db)

	// Services
	authenticateService := authSrv.NewService(userRepo, authenticator, time.Hour*24)
	articleService := article.NewService(articleRepo, &validation.Article{})
	categoryService := category.NewService(categoryRepo, &validation.Category{})
	roleService := role.NewService(roleRepo, &validation.Role{})
	userService := user.NewService(userRepo, &validation.User{})

	services := httpBroker.Services{
		Auth:     authenticateService,
		Article:  articleService,
		Category: categoryService,
		Role:     roleService,
		User:     userService,
	}

	return &services
}
