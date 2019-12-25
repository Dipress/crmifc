package schema

import "database/sql"

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

const seeds = `
	-- Create role with name "Admin"
	INSERT INTO roles (name) VALUES('Admin') ON CONFLICT DO NOTHING;
	INSERT INTO roles (name) VALUES('Manager') ON CONFLICT DO NOTHING;

	-- Create admin with password "password123"
	INSERT INTO users (role_id, username, email, password_hash) VALUES 
	(1, 'Admin', 'admin@example.com', '$2a$10$lGMGO59qq7yKx.zwtI4cZul5lM7YVS1v07.4hlSAPrbngUDfddQBK') 
	ON CONFLICT DO NOTHING;

	-- Create manager with password "password123"
	INSERT INTO users (role_id, username, email, password_hash) VALUES 
	(2, 'Manager', 'manager@example.com', '$2a$10$lGMGO59qq7yKx.zwtI4cZul5lM7YVS1v07.4hlSAPrbngUDfddQBK') 
	ON CONFLICT DO NOTHING;
`
