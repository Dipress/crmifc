package update

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	gomock "github.com/golang/mock/gomock"
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
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, nil)
				m.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
			name: "find role error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "update role error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, nil)
				m.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock errpr"))
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

			form := role.Form{
				Name: "Manager",
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
