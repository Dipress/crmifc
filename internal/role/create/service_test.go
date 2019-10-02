package create

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Create(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, f *role.NewRole, rol *role.Role) error
		validaterFunc  func(ctx context.Context, form *role.Form) error
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, f *role.NewRole, rol *role.Role) error {
				return nil
			},
			validaterFunc: func(ctx context.Context, form *role.Form) error {
				return nil
			},
		},
		{
			name: "validation",
			repositoryFunc: func(ctx context.Context, f *role.NewRole, rol *role.Role) error {
				return nil
			},
			validaterFunc: func(ctx context.Context, form *role.Form) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "create role",
			repositoryFunc: func(ctx context.Context, f *role.NewRole, rol *role.Role) error {
				return errors.New("mock error")
			},
			validaterFunc: func(ctx context.Context, form *role.Form) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(repositoryFunc(tc.repositoryFunc), validaterFunc(tc.validaterFunc))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			form := role.Form{
				Name: "Admin",
			}

			_, err := s.Create(ctx, &form)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, f *role.NewRole, rol *role.Role) error

func (r repositoryFunc) CreateRole(ctx context.Context, f *role.NewRole, rol *role.Role) error {
	return r(ctx, f, rol)
}

type validaterFunc func(ctx context.Context, form *role.Form) error

func (v validaterFunc) Validate(ctx context.Context, form *role.Form) error {
	return v(ctx, form)
}
