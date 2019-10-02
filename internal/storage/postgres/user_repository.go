package postgres

import (
	"context"
	"database/sql"

	"github.com/dipress/crmifc/internal/auth"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	driverName = "postgres"
)

// Repository holds CRUD actions.
type Repository struct {
	db *sqlx.DB
}

//NewRepository factory prepares the repository to work.
func NewRepository(db *sql.DB) *Repository {
	r := Repository{
		db: sqlx.NewDb(db, driverName),
	}

	return &r
}

const createUserQuery = `INSERT INTO 
	users (username, email, password_hash, role_id) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, role_id, username, email, created_at, updated_at`

// CreateUser insert a new user into the database.
func (r *Repository) CreateUser(ctx context.Context, f *user.NewUser, usr *user.User) error {
	if err := r.db.QueryRowContext(ctx, createUserQuery, f.Username, f.Email, f.PasswordHash, f.RoleID).
		Scan(&usr.ID, &usr.Role.ID, &usr.Username, &usr.Email, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
		return errors.Wrap(err, "query context scan")
	}
	return nil
}

const findUserQuery = `SELECT id, username, email, created_at, updated_at, role_id FROM users WHERE id = $1`

// FindUser finds a user by id.
func (r *Repository) FindUser(ctx context.Context, id int) (*user.User, error) {
	var u user.User
	if err := r.db.QueryRowContext(ctx, findUserQuery, id).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Role.ID,
		); err != nil {
		if err == sql.ErrNoRows {
			return nil, role.ErrNotFound
		}
		return nil, errors.Wrap(err, "query row scan")
	}
	return &u, nil
}

const updateUserQuery = `UPDATE users SET username=:username, email=:email, password_hash=:password_hash, role_id=:role_id, updated_at=now() WHERE id=:id`

// UpdateUser updates user by id.
func (r *Repository) UpdateUser(ctx context.Context, id int, u *user.User) error {
	stmt, err := r.db.PrepareNamed(updateUserQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id":            id,
		"username":      u.Username,
		"email":         u.Email,
		"password_hash": u.PasswordHash,
		"role_id":       u.Role.ID,
	}); err != nil {
		if err == sql.ErrNoRows {
			return user.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}
	return nil
}

const deleteUserQuery = `DELETE FROM users WHERE id=:id`

// DeleteUser deletes user by id.
func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	stmt, err := r.db.PrepareNamed(deleteUserQuery)
	if err != nil {
		return errors.Wrap(err, "prepare named")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		if err == sql.ErrNoRows {
			return user.ErrNotFound
		}
		return errors.Wrap(err, "exec context")
	}
	return nil
}

const uniqueUsernameQuery = `SELECT COUNT(*) FROM users WHERE username = $1`

// UniqueUsername checks that username is unique.
func (r *Repository) UniqueUsername(ctx context.Context, username string) error {
	var c int
	if err := r.db.QueryRowContext(ctx, uniqueUsernameQuery, username).Scan(&c); err != nil {
		return errors.Wrap(err, "scan error")
	}

	if c > 0 {
		return user.ErrUsernameExists
	}

	return nil
}

const uniqueEmailQuery = `SELECT COUNT(*) FROM users WHERE email = $1`

// UniqueEmail checks that email address is unique.
func (r *Repository) UniqueEmail(ctx context.Context, email string) error {
	var c int
	if err := r.db.QueryRowContext(ctx, uniqueEmailQuery, email).Scan(&c); err != nil {
		return errors.Wrap(err, "scan error")
	}

	if c > 0 {
		return user.ErrEmailExists
	}

	return nil
}

const emailFindQuery = `SELECT 
	id, username, email, password_hash, created_at, updated_at 
	FROM users WHERE email = $1`

// FindByEmail finds users by e-mail.
func (r *Repository) FindByEmail(ctx context.Context, email string, usr *user.User) error {
	if err := r.db.QueryRowContext(ctx, emailFindQuery, email).
		Scan(
			&usr.ID,
			&usr.Username,
			&usr.Email,
			&usr.PasswordHash,
			&usr.CreatedAt,
			&usr.UpdatedAt,
		); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrEmailNotFound
		}

		return errors.Wrap(err, "scan error")
	}

	return nil
}

const listUsersQuery = `SELECT 
	users.id, users.email, users.email, users.created_at, users.updated_at, roles.id, roles.name 
	FROM users
	LEFT JOIN roles ON users.role_id = roles.id`

// ListUsers returns all users.
func (r *Repository) ListUsers(ctx context.Context, usr *user.Users) error {
	rows, err := r.db.QueryxContext(ctx, listUsersQuery)
	if err != nil {
		return errors.Wrap(err, "query rows")
	}
	defer rows.Close()

	for rows.Next() {
		var user user.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Role.ID,
			&user.Role.Name,
		); err != nil {
			return errors.Wrap(err, "users query row scan on loop")
		}

		usr.Users = append(usr.Users, user)
	}

	return nil
}