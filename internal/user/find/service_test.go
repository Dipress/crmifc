package find

import (
	"context"
	"errors"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, id int) (*user.User, error)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, id int) (*user.User, error) {
				return &user.User{}, nil
			},
		},
		{
			name: "repository error",
			repositoryFunc: func(ctx context.Context, id int) (*user.User, error) {
				return &user.User{}, errors.New("mock error")
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

			_, err := s.FindUser(ctx, 1)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, id int) (*user.User, error)

func (r repositoryFunc) FindUser(ctx context.Context, id int) (*user.User, error) {
	return r(ctx, id)
}
