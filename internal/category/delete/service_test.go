package delete

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/category"
	gomock "github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Delete(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, nil)
				m.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find category error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "delete category error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, nil)
				m.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)
			tc.repositoryFunc(repo)

			s := NewService(repo)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := s.Delete(ctx, 1)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
