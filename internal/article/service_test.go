package article

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/kit/auth"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Service_Create(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		validaterFunc  func(mock *MockValidater)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:           "validation error",
			repositoryFunc: func(m *MockRepository) {},
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "create article error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
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
			validater := NewMockValidater(ctrl)

			tc.repositoryFunc(repo)
			tc.validaterFunc(validater)

			s := NewService(repo, validater)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			claims := auth.Claims{}
			newCtx := auth.ToContext(ctx, &claims)

			form := Form{
				CategoryID: 1,
				Title:      "my title",
				Body:       "my body",
			}
			_, err := s.Create(newCtx, &form)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_Service_Find(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, nil)
			},
		},
		{
			name: "internal error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, errors.New("mock error"))
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

			s := NewService(repo, nil)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			claims := auth.Claims{}
			newCtx := auth.ToContext(ctx, &claims)

			_, err := s.Find(newCtx, 1)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func Test_Service_Update(t *testing.T) {
	tests := []struct {
		name           string
		validaterFunc  func(mock *MockValidater)
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
			name: "find article error",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "update article error",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
			validater := NewMockValidater(ctrl)

			tc.repositoryFunc(repo)
			tc.validaterFunc(validater)

			s := NewService(repo, validater)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			claims := auth.Claims{}
			newCtx := auth.ToContext(ctx, &claims)

			form := Form{
				CategoryID: 15,
				Title:      "update my awesome titie",
				Body:       "update my awesome body",
			}

			_, err := s.Update(newCtx, 1, &form)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_Service_Delete(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, nil)
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find article error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "delete article error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Article{}, nil)
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("mock errror"))
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

			s := NewService(repo, nil)

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
