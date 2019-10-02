package postgres

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/dipress/crmifc/internal/kit/docker"
	"github.com/dipress/crmifc/internal/storage/postgres/schema"
	"github.com/ory/dockertest"
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

	txdb.Register("pgsqltx", "postgres", fmt.Sprintf("password=test user=test dbname=test host=localhost port=%s sslmode=disable", pgDocker.Resource.GetPort("5432/tcp")))

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
