package update

import (
	"context"
	"testing"

	"github.com/dipress/crmifc/internal/category"
	gomock "github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Update(t *testing.T) {
	tests := []struct {
		name           string
		validateFunc   func(mock *MockValidater)
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, nil)
				m.EXPECT().UpdateCategory(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "validation error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			repositoryFunc: func(m *MockRepository) {},
			wantErr:        true,
		},
		{
			name: "find category error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "update category error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindCategory(gomock.Any(), gomock.Any()).Return(&category.Category{}, nil)
				m.EXPECT().UpdateCategory(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock errpr"))
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			validater := NewMockValidater(ctrl)
			repository := NewMockRepository(ctrl)
			tc.validateFunc(validater)
			tc.repositoryFunc(repository)

			s := NewService(repository, validater)

			form := category.Form{
				Name: "Region Channels",
			}

			_, err := s.Update(ctx, 1, &form)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
