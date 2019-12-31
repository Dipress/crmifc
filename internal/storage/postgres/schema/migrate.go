package schema

import (
	"database/sql"
	"fmt"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
)

const (
	source          = "go-bindata"
	database        = "postgres"
	migrationsTable = "versions"
)

// go:generate go-bindata -prefix migrations/ -pkg schema -o migrations.bindata.go migrations/

// Migrate migrates schema to given database connection.
func Migrate(db *sql.DB) error {
	m, err := newMigration(db)
	if err != nil {
		return fmt.Errorf("new migration: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}

	return nil
}

func newMigration(db *sql.DB) (*migrate.Migrate, error) {
	r := bindata.Resource(AssetNames(), Asset)
	s, err := bindata.WithInstance(r)
	if err != nil {
		return nil, fmt.Errorf("prepare source instance: %w", err)
	}

	cfg := postgres.Config{
		MigrationsTable: migrationsTable,
	}

	d, err := postgres.WithInstance(db, &cfg)
	if err != nil {
		return nil, fmt.Errorf("prepare database instance: %w", err)
	}

	m, err := migrate.NewWithInstance(source, s, database, d)
	if err != nil {
		return nil, fmt.Errorf("prepare migrate instance: %w", err)
	}

	return m, nil
}
