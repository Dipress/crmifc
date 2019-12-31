package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dipress/crmifc/internal/role"
	"github.com/jmoiron/sqlx"
)

// RoleRepository holds CRUD actions.
type RoleRepository struct {
	db *sqlx.DB
}

//NewRoleRepository factory prepares the repository to work.
func NewRoleRepository(db *sql.DB) *RoleRepository {
	r := RoleRepository{
		db: sqlx.NewDb(db, driverName),
	}

	return &r
}

const createRoleQuery = `INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at, updated_at`

// Create insert a new role into the database.
func (r *RoleRepository) Create(ctx context.Context, f *role.NewRole, rol *role.Role) error {
	if err := r.db.QueryRowContext(ctx, createRoleQuery, f.Name).
		Scan(&rol.ID, &rol.Name, &rol.CreatedAt, &rol.UpdatedAt); err != nil {
		return fmt.Errorf("query context scan: %w", err)
	}
	return nil
}

const findRoleQuery = `SELECT id, name, created_at, updated_at FROM roles where id = $1`

// Find finds a role by id.
func (r *RoleRepository) Find(ctx context.Context, id int) (*role.Role, error) {
	var rol role.Role
	if err := r.db.QueryRowContext(ctx, findRoleQuery, id).
		Scan(&rol.ID, &rol.Name, &rol.CreatedAt, &rol.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, role.ErrNotFound
		}
		return nil, fmt.Errorf("query row scan: %w", err)
	}
	return &rol, nil
}

const updateRoleQuery = `UPDATE roles SET name=:name, updated_at=now() WHERE id=:id`

// Update updates role by id.
func (r *RoleRepository) Update(ctx context.Context, id int, rl *role.Role) error {

	stmt, err := r.db.PrepareNamed(updateRoleQuery)
	if err != nil {
		return fmt.Errorf("prepare named: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id":   id,
		"name": rl.Name,
	}); err != nil {
		if err == sql.ErrNoRows {
			return role.ErrNotFound
		}
		return fmt.Errorf("exec context: %w", err)
	}
	return nil
}

const deleteRoleQuery = `DELETE FROM roles WHERE id=:id`

// Delete deletes role by id.
func (r *RoleRepository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareNamed(deleteRoleQuery)
	if err != nil {
		return fmt.Errorf("prepare named: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		if err == sql.ErrNoRows {
			return role.ErrNotFound
		}
		return fmt.Errorf("exec context: %w", err)
	}
	return nil
}

const listRoleQuery = `SELECT * FROM roles`

// List shows all roles.
func (r *RoleRepository) List(ctx context.Context, roles *role.Roles) error {
	rows, err := r.db.QueryxContext(ctx, listRoleQuery)
	if err != nil {
		return fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rl role.Role
		if err := rows.Scan(&rl.ID, &rl.Name, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
			return fmt.Errorf("roles query row scan on loop: %w", err)
		}

		roles.Roles = append(roles.Roles, rl)
	}

	return nil
}
