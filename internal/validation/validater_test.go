package validation

import (
	"context"
	"reflect"
	"testing"

	"github.com/dipress/crmifc/internal/category"
	"github.com/dipress/crmifc/internal/role"
	"github.com/dipress/crmifc/internal/user"
)

func TestRoleValidate(t *testing.T) {
	tests := []struct {
		name    string
		form    role.Form
		wantErr bool
		expect  Errors
	}{
		{
			name: "ok",
			form: role.Form{
				Name: "Admin",
			},
		},
		{
			name:    "blank name",
			form:    role.Form{},
			wantErr: true,
			expect: Errors{
				"name": "cannot be blank",
			},
		},
		{
			name: "long name",
			form: role.Form{
				Name: "This is long name for role, this title is way larger is allowed one",
			},
			wantErr: true,
			expect: Errors{
				"name": "the length must be between 1 and 50",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var r Role
			err := r.Validate(ctx, &tc.form)
			if tc.wantErr {
				got, ok := err.(Errors)
				if !ok {
					t.Errorf("unknown error: %v", err)
					return
				}

				if !reflect.DeepEqual(tc.expect, got) {
					t.Errorf("expected: %+#v got: %+#v", tc.expect, got)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})

	}
}

func TestUserValidate(t *testing.T) {
	tests := []struct {
		name    string
		form    user.Form
		wantErr bool
		expect  Errors
	}{
		{
			name: "ok",
			form: user.Form{
				Username: "Shepard",
				Email:    "shepard@normandy.com",
				Password: "mypassword",
				RoleID:   1,
			},
		},
		{
			name: "blank username",
			form: user.Form{
				Email:    "shepard@normandy.com",
				Password: "mypassword",
				RoleID:   1,
			},
			wantErr: true,
			expect: Errors{
				"username": "cannot be blank",
			},
		},
		{
			name: "long username",
			form: user.Form{
				Username: "This is long username for user, this title is way larger is allowed one",
				Email:    "shepard@normandy.com",
				Password: "mypassword",
				RoleID:   1,
			},
			wantErr: true,
			expect: Errors{
				"username": "the length must be between 1 and 50",
			},
		},
		{
			name: "blank email",
			form: user.Form{
				Username: "Shepard",
				Password: "mypassword",
				RoleID:   1,
			},
			wantErr: true,
			expect: Errors{
				"email": "cannot be blank",
			},
		},
		{
			name: "not valid email",
			form: user.Form{
				Email:    "shepardnormandy.com",
				Username: "Shepard",
				Password: "mypassword",
				RoleID:   1,
			},
			wantErr: true,
			expect: Errors{
				"email": "must be a valid email address",
			},
		},
		{
			name: "blank role_id",
			form: user.Form{
				Username: "Shepard",
				Email:    "shepard@normandy.com",
				Password: "mypassword",
			},
			wantErr: true,
			expect: Errors{
				"role_id": "cannot be blank",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var u User
			err := u.Validate(ctx, &tc.form)
			if tc.wantErr {
				got, ok := err.(Errors)
				if !ok {
					t.Errorf("unknown error: %v", err)
					return
				}

				if !reflect.DeepEqual(tc.expect, got) {
					t.Errorf("expected: %+#v got: %+#v", tc.expect, got)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCategoryValidate(t *testing.T) {
	tests := []struct {
		name    string
		form    category.Form
		wantErr bool
		expect  Errors
	}{
		{
			name: "ok",
			form: category.Form{
				Name: "Real IPs",
			},
		},
		{
			name:    "blank name",
			form:    category.Form{},
			wantErr: true,
			expect: Errors{
				"name": "cannot be blank",
			},
		},
		{
			name: "long name",
			form: category.Form{
				Name: "This is long name for category, this title is way larger is allowed one.",
			},
			wantErr: true,
			expect: Errors{
				"name": "the length must be between 1 and 50",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var c Category
			err := c.Validate(ctx, &tc.form)
			if tc.wantErr {
				got, ok := err.(Errors)
				if !ok {
					t.Errorf("unknown error: %v", err)
					return
				}

				if !reflect.DeepEqual(tc.expect, got) {
					t.Errorf("expected: %+#v got: %+#v", tc.expect, got)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})

	}
}
