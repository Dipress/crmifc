package main

import (
	"crypto/rsa"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
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
	db            *sql.DB
	authenticator *auth.Authenticator
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
	keyContents, err := ioutil.ReadFile("../internal/kit/keys/jwtRS256.key")
	if err != nil {
		log.Fatalf("reading auth private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
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
