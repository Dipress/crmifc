package find

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/category"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Find(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(ctx context.Context, id int) (*category.Category, error)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(ctx context.Context, id int) (*category.Category, error) {
				return &category.Category{}, nil
			},
		},
		{
			name: "repository error",
			repositoryFunc: func(ctx context.Context, id int) (*category.Category, error) {
				return &category.Category{}, errors.New("mock error")
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
			_, err := s.Repository.FindCategory(ctx, 1)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

type repositoryFunc func(ctx context.Context, id int) (*category.Category, error)

func (r repositoryFunc) FindCategory(ctx context.Context, id int) (*category.Category, error) {
	return r(ctx, id)
}
