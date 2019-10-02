package update

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Update(t *testing.T) {
	tests := []struct {
		name           string
		validaterFunc  func(m *MockValidater)
		repositoryFunc func(m *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, nil)
				m.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "validation error",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			repositoryFunc: func(m *MockRepository) {},
			wantErr:        true,
		},
		{
			name: "find user error",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "update user error",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(&user.User{}, nil)
				m.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
			repo := NewMockRepository(ctrl)

			tc.validaterFunc(validater)
			tc.repositoryFunc(repo)

			s := NewService(repo, validater)

			form := user.Form{
				Username: "Dipress",
				Email:    "dipress@example.com",
				Password: "myNewPassword",
				RoleID:   6,
			}

			err := s.Update(ctx, 5, &form)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
