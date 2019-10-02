package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/user"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(m *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, nil)
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find user error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "delete user error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, nil)
				m.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
