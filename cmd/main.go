package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	httpBroker "github.com/dipress/crmifc/internal/broker/http"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/storage/postgres/schema"
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
		privateKeyFile = flag.String("key", "../internal/kit/keys/jwtRS256.key", "private key file path")
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
	authenticator, err := auth.NewAuthenticator(key, *keyID, alg, publicKeyLookup)
	if err != nil {
		log.Fatalf("constructing authenticator: %v", err)
	}

	// Setup handlers.
	srv := setupServer(*addr, db, authenticator)
	if err := srv.ListenAndServe(); err != nil {
		errors.Wrap(err, "filed to serve http")
	}
}

func setupServer(addr string, db *sql.DB, authenticator *auth.Authenticator) *http.Server {
	return httpBroker.NewServer(addr, db, authenticator)
}
