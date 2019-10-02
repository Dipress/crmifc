package create

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
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
				m.EXPECT().UniqueUsername(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().UniqueEmail(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
			name: "username unique error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().UniqueUsername(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "email unique error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().UniqueUsername(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().UniqueEmail(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "create error",
			validateFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().UniqueUsername(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().UniqueEmail(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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

			validator := NewMockValidater(ctrl)
			repo := NewMockRepository(ctrl)

			tc.validateFunc(validator)
			tc.repositoryFunc(repo)

			s := NewService(repo, validator)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			form := user.Form{
				Username: "username",
				Email:    "username@example.com",
				Password: "password123",
				RoleID:   1,
			}

			var u user.User
			err := s.Create(ctx, &form, &u)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
