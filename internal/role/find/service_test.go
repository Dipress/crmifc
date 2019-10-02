package find

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Find(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, id int) (*role.Role, error)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, id int) (*role.Role, error) {
				return &role.Role{}, nil
			},
		},
		{
			name: "respository error",
			repositoryFunc: func(ctx context.Context, id int) (*role.Role, error) {
				return &role.Role{}, errors.New("mock error")
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := NewService(repositoryFunc(tc.repositoryFunc))

			_, err := s.FindRole(ctx, 1)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, id int) (*role.Role, error)

func (r repositoryFunc) FindRole(ctx context.Context, id int) (*role.Role, error) {
	return r(ctx, id)
}
