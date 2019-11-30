package role

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Create_Service(t *testing.T) {
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
			name: "create role error",
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

			form := Form{
				Name: "Admin",
			}

			_, err := s.Create(ctx, &form)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func Test_Find_Service(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, nil)
			},
		},
		{
			name: "internal error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, errors.New("mock error"))
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

			_, err := s.Find(ctx, 1)
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
		repositoryFunc func(mock *MockRepository)
		validaterFunc  func(mock *MockValidater)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
			name: "find role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, errors.New("mock error"))
			},
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "update role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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

			form := Form{
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

func Test_Service_Delete(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, nil)
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "find role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "delete role error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&Role{}, nil)
				m.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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

func Test_List_Service(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "internal error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().List(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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

			_, err := s.List(ctx)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
