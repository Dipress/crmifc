package update

import (
	"context"
	"errors"
	"testing"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/kit/auth"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().FindArticle(gomock.Any(), gomock.Any()).Return(&article.Article{}, nil)
				m.EXPECT().UpdateArticle(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "validation",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			repositoryFunc: func(m *MockRepository) {},
			wantErr:        true,
		},
		{
			name: "find user",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "find article",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().FindArticle(gomock.Any(), gomock.Any()).Return(&article.Article{}, errors.New("mock error"))
			},
			wantErr: true,
		},
		{
			name: "update article",
			validaterFunc: func(m *MockValidater) {
				m.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
			},
			repositoryFunc: func(m *MockRepository) {
				m.EXPECT().FindByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().FindArticle(gomock.Any(), gomock.Any()).Return(&article.Article{}, nil)
				m.EXPECT().UpdateArticle(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
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
