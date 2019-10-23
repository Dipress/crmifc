package create

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/kit/auth"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	tests := []struct {
		name           string
		repositoryFunc func(mock *MockRepository)
		validaterFunc  func(mock *MockValidater)
		wantErr        bool
	}{
		{
			name: "ok",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().CreateArticle(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
			name: "find user error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "create article error",
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().CreateArticle(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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

			form := article.Form{
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
