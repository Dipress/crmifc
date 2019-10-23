package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/role"
	"github.com/golang/mock/gomock"
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
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, nil)
				m.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "delete role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindRole(gomock.Any(), gomock.Any()).Return(&role.Role{}, nil)
				m.EXPECT().DeleteRole(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
