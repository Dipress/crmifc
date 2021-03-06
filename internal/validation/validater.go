package validation

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/user"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	mismatchMsg   = "mismatch"
	validationMsg = "you have validation errors"
)

// Errors holds validation errors.
type Errors map[string]string

// Error implements error interface.
func (v Errors) Error() string {
	return validationMsg
}

// Role holds form validations.
type Role struct{}

// Validate validates role form.
func (r *Role) Validate(ctx context.Context, form *role.Form) error {
	ves := make(Errors)

	if err := validation.Validate(form.Name,
		validation.Required,
		validation.Length(1, 50)); err != nil {
		ves["name"] = err.Error()
	}

	if len(ves) > 0 {
		return ves
	}

	return nil
}

// User holds form validations.
type User struct{}

// Validate validates user form.
func (u *User) Validate(ctx context.Context, form *user.Form) error {
	ves := make(Errors)

	if err := validation.Validate(form.Username,
		validation.Required,
		validation.Length(1, 50)); err != nil {
		ves["username"] = err.Error()
	}

	if err := validation.Validate(form.Email,
		validation.Required,
		is.Email,
		validation.Length(1, 50)); err != nil {
		ves["email"] = err.Error()
	}

	if err := validation.Validate(form.Password,
		validation.Required,
		validation.Length(1, 72)); err != nil {
		ves["password"] = err.Error()
	}

	if err := validation.Validate(form.RoleID,
		validation.Required); err != nil {
		ves["role_id"] = err.Error()
	}

	if len(ves) > 0 {
		return ves
	}
	return nil
}

// Category holds form validations.
type Category struct{}

// Validate validates category form.
func (c *Category) Validate(ctx context.Context, form *category.Form) error {
	ves := make(Errors)

	if err := validation.Validate(form.Name,
		validation.Required,
		validation.Length(1, 50)); err != nil {
		ves["name"] = err.Error()
	}

	if len(ves) > 0 {
		return ves
	}

	return nil
}

// Article holds form validations.
type Article struct{}

// Validate validates article form.
func (a *Article) Validate(ctx context.Context, form *article.Form) error {
	ves := make(Errors)

	if err := validation.Validate(form.CategoryID,
		validation.Required); err != nil {
		ves["category_id"] = err.Error()
	}

	if err := validation.Validate(form.Title,
		validation.Required); err != nil {
		ves["title"] = err.Error()
	}

	if err := validation.Validate(form.Body,
		validation.Required); err != nil {
		ves["body"] = err.Error()
	}

	if len(ves) > 0 {
		return ves
	}
	return nil
}
