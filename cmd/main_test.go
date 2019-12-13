package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/dgrijalva/jwt-go"
	"github.com/dipress/crmifc/internal/kit/auth"
	"github.com/dipress/crmifc/internal/kit/docker"
	"github.com/dipress/crmifc/internal/storage/postgres"
	"github.com/dipress/crmifc/internal/storage/postgres/schema"
	"github.com/ory/dockertest"

	_ "github.com/lib/pq"
)

const (
	caseTimeout = 5 * time.Second
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(m.Run())
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	pgDocker, err := docker.NewPostgres(pool)
	if err != nil {
		log.Fatalf("prepare postgres with docker: %v", err)
	}
	db = pgDocker.DB

	if err := schema.Migrate(db); err != nil {
		log.Fatalf("migrate schema: %v", err)
	}

	txdb.Register("pgsqltx", "postgres",
		fmt.Sprintf("password=test user=test dbname=test host=localhost port=%s sslmode=disable",
			pgDocker.Resource.GetPort("5432/tcp")),
	)

	code := m.Run()

	db.Close()
	if err := pool.Purge(pgDocker.Resource); err != nil {
		log.Fatalf("could not purge postgres docker: %v", err)
	}

	os.Exit(code)
}

func postgresDB(t *testing.T) (db *sql.DB, teardown func() error) {
	dbName := fmt.Sprintf("db_%d", time.Now().UnixNano())
	db, err := sql.Open("pgsqltx", dbName)

	if err != nil {
		t.Fatalf("open postgres tx connection: %s", err)
	}

	return db, db.Close
}

func authenticatorSetup(db *sql.DB) *auth.Authenticator {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtKey))
	if err != nil {
		log.Fatalf("parsing auth private key: %v", err)
	}

	publicKeyLookup := auth.NewSingleKeyFunc("12345", key.Public().(*rsa.PublicKey))

	userRepo := postgres.NewUserRepository(db)
	ac, err := auth.NewAuthenticator(key, "12345", alg, publicKeyLookup, userRepo)
	if err != nil {
		log.Fatalf("constructing authenticator: %v", err)
	}

	return ac
}

var jwtKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
