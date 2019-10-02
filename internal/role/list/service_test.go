package list

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Service_List(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, roles *role.Roles) error
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, roles *role.Roles) error {
				return nil
			},
		},
		{
			name: "repository error",
			repositoryFunc: func(ctx context.Context, roles *role.Roles) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(repositoryFunc(tc.repositoryFunc))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err := s.List(ctx)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, roles *role.Roles) error

func (r repositoryFunc) ListRoles(ctx context.Context, roles *role.Roles) error {
	return r(ctx, roles)
}
