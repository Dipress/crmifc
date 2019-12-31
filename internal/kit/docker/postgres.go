package docker

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

const (
	dockerStartWait = 30 * time.Second
)

// PostgresDocker holds connection
// to the db and resource info
// for shutdown.
type PostgresDocker struct {
	DB       *sql.DB
	Resource *dockertest.Resource
}

// NewPostgres starts lates postgres docker image
// and tries to connect with it.
func NewPostgres(pool *dockertest.Pool) (*PostgresDocker, error) {
	res, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_PASSWORD=test",
		"POSTGRES_USER=test",
		"POSTGRES_DB=test",
	})
	if err != nil {
		return nil, fmt.Errorf("start postgres: %w", err)
	}

	purge := func() {
		pool.Purge(res)
	}

	errChan := make(chan error)
	done := make(chan struct{})

	var db *sql.DB

	go func() {
		if err := pool.Retry(func() error {
			db, err = sql.Open("postgres", fmt.Sprintf("user=test password=test dbname=test host=localhost port=%s sslmode=disable", res.GetPort("5432/tcp")))
			if err != nil {
				return err
			}
			if err := db.Ping(); err != nil {
				db.Close()
				return err
			}
			return nil
		}); err != nil {
			errChan <- err
		}

		close(done)
	}()

	select {
	case err := <-errChan:
		purge()
		return nil, fmt.Errorf("check connection: %w", err)
	case <-time.After(dockerStartWait):
		purge()
		return nil, errors.New("timeout on checking postgres connection")
	case <-done:
		close(errChan)
	}

	pd := PostgresDocker{
		DB:       db,
		Resource: res,
	}

	return &pd, nil
}
