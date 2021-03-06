package postgres

import (
	"context"
	"database/sql"

	"github.com/dipress/crmifc/internal/role"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "query context scan")
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
		return nil, errors.Wrap(err, "query row scan")
	}
	return &rol, nil
}

const updateRoleQuery = `UPDATE roles SET name=:name, updated_at=now() WHERE id=:id`

// Update updates role by id.
func (r *RoleRepository) Update(ctx context.Context, id int, rl *role.Role) error {

	stmt, err := r.db.PrepareNamed(updateRoleQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id":   id,
		"name": rl.Name,
	}); err != nil {
		if err == sql.ErrNoRows {
			return role.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}
	return nil
}

const deleteRoleQuery = `DELETE FROM roles WHERE id=:id`

// Delete deletes role by id.
func (r *RoleRepository) Delete(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareNamed(deleteRoleQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		if err == sql.ErrNoRows {
			return role.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}
	return nil
}

const listRoleQuery = `SELECT * FROM roles`

// List shows all roles.
func (r *RoleRepository) List(ctx context.Context, roles *role.Roles) error {
	rows, err := r.db.QueryxContext(ctx, listRoleQuery)
	if err != nil {
		return errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	for rows.Next() {
		var rl role.Role
		if err := rows.Scan(&rl.ID, &rl.Name, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
			return errors.Wrap(err, "roles query row scan on loop")
		}

		roles.Roles = append(roles.Roles, rl)
	}

	return nil
}
