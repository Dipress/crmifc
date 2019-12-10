package article

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"github.com/dipress/crmifc/internal/article"
	"github.com/dipress/crmifc/internal/validation"
)

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name        string
		serviceFunc func(m *MockService)
		code        int
	}{
		{
			name: "ok",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&article.Article{}, nil)
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&article.Article{}, make(validation.Errors))
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&article.Article{}, errors.New("mock error"))
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			tc.serviceFunc(service)

			h := CreateHandler{service}
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

func TestFindHandler(t *testing.T) {
	tests := []struct {
		name        string
		serviceFunc func(m *MockService)
		code        int
	}{
		{
			name: "ok",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&article.Article{}, nil)
			},
			code: http.StatusOK,
		},
		{
			name: "internal error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Find(gomock.Any(), gomock.Any()).Return(&article.Article{}, errors.New("mock error"))
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			tc.serviceFunc(service)

			h := FindHandler{service}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name        string
		serviceFunc func(m *MockService)
		code        int
	}{
		{
			name: "ok",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&article.Article{}, nil)
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&article.Article{}, make(validation.Errors))
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&article.Article{}, errors.New("mock error"))
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			tc.serviceFunc(service)

			h := UpdateHandler{service}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name        string
		serviceFunc func(m *MockService)
		code        int
	}{
		{
			name: "ok",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			code: http.StatusOK,
		},
		{
			name: "repository error",
			serviceFunc: func(mock *MockService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			tc.serviceFunc(service)

			h := DeleteHandler{service}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "http://example.com", strings.NewReader("{}"))
			r = mux.SetURLVars(r, map[string]string{"id": "1"})

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}

func TestListHandler(t *testing.T) {
	tests := []struct {
		name        string
		serviceFunc func(mock *MockService)
		code        int
	}{
		{
			name: "ok",
			serviceFunc: func(m *MockService) {
				m.EXPECT().List(gomock.Any()).Return(&article.Articles{}, nil)
			},
			code: http.StatusOK,
		},
		{
			name: "repository error",
			serviceFunc: func(m *MockService) {
				m.EXPECT().List(gomock.Any()).Return(&article.Articles{}, errors.New("mock error"))
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			tc.serviceFunc(service)

			h := ListHandler{service}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://example.com", strings.NewReader("{}"))

			err := h.Handle(w, r)
			if w.Code != tc.code {
				t.Errorf("unexpected code: %d expected %d error: %v", w.Code, tc.code, err)
			}
		})
	}
}
