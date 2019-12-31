package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/dgrijalva/jwt-go"
	"github.com/dipress/crmifc/internal/article"
	authSrv "github.com/dipress/crmifc/internal/auth"
	httpBroker "github.com/dipress/crmifc/internal/broker/http"
	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/kit/logger"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/storage/postgres/schema"
	"github.com/dipress/crmifc/internal/user"
	"github.com/dipress/crmifc/internal/validation"
	"github.com/mattes/migrate"
)

const (
	alg = "RS256"
)

func main() {
	var (
		addr           = flag.String("addr", ":8080", "address of http server")
		dsn            = flag.String("dsn", "", "postgres database DSN")
		privateKeyFile = flag.String("key", "", "private key file path")
		keyID          = flag.String("id", "123456", "private key id")
	)
	flag.Parse()

	// Logger initialize.
	logger, err := logger.New(
		logger.WithSensitiveFields([]string{
			"password",
			"token",
		}),
	)

	if err != nil {
		log.Fatalf("prepare logger: %v\n", err)
	}

	// Setup database connection.
	logger.Info("connecting to db", nil)
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create db: %w", err), nil)
	}

	if err := db.Ping(); err != nil {
		logger.Fatal(fmt.Errorf("failed to connect db: %w", err), nil)
	}

	defer db.Close()
	logger.Info("connection to db established", nil)

	// Migrate schema.
	if err := schema.Migrate(db); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Fatal(fmt.Errorf("failed to migrate schema: %w", err), nil)
		}
	}

	// Seed data.
	if err := schema.Seed(db); err != nil {
		logger.Fatal(fmt.Errorf("failed to seeds data: %w", err), nil)
	}

	// Reppositories.
	userRepo := postgres.NewUserRepository(db)

	// Authentication setup.
	keyContents, err := ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		logger.Fatal(fmt.Errorf("reading auth private key: %w", err), nil)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		logger.Fatal(fmt.Errorf("parsing auth private key: %w", err), nil)
	}

	publicKeyLookup := auth.NewSingleKeyFunc(*keyID, key.Public().(*rsa.PublicKey))
	authenticator, err := auth.NewAuthenticator(key, *keyID, alg, publicKeyLookup, userRepo)
	if err != nil {
		logger.Fatal(fmt.Errorf("constructing authenticator: %w", err), nil)
	}

	// Make a channel for errors.
	errChan := make(chan error)

	// Services
	services := setupServices(db, authenticator)

	// Setup server.
	srv := setupServer(*addr, logger, services, authenticator)

	go func() {
		logger.Info(fmt.Sprintf("starting %s server", srv.Addr), nil)
		if err := srv.ListenAndServe(); err != nil {
			errChan <- fmt.Errorf("launch server %s: %w", srv.Addr, err)
		}
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Fatal(err, nil)
	case <-osSignals:
		if err := srv.Shutdown(context.TODO()); err != nil {
			logger.Fatal(fmt.Errorf("stop server %s: %w", srv.Addr, err), nil)
		}
	}
}

func setupServer(addr string, logger *logger.Logger, services *httpBroker.Services, authenticator *auth.Authenticator) *http.Server {
	return httpBroker.NewServer(addr, logger, services, authenticator)
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
