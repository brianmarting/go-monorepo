package handlers

import (
	"errors"
	"go-monorepo/identity-provider/service/mocks"
	"go-monorepo/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var user = model.User{
	Id:           32,
	ExternalId:   "dq9b5f7e-61e6-480e-91bb-269f7a929622",
	Name:         "John",
	Password:     "somepwd",
	TokenVersion: 1,
}

func Test_userHandler_CreateUser(t *testing.T) {
	serviceMock := new(mocks.UserServiceMock)
	handler := userHandler{
		service: serviceMock,
	}

	tests := []struct {
		name         string
		rec          *httptest.ResponseRecorder
		reqFn        func() *http.Request
		mockFn       func() *mock.Call
		expectResult bool
		expectError  bool
	}{
		{
			name: "Should create user",
			rec:  httptest.NewRecorder(),
			reqFn: func() *http.Request {
				return httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			},
			mockFn: func() *mock.Call {
				return serviceMock.
					On("Create", mock.Anything).
					Return(nil)
			},
			expectResult: true,
		},
		{
			name: "Should return err on bad body",
			rec:  httptest.NewRecorder(),
			reqFn: func() *http.Request {
				return httptest.NewRequest("POST", "/", nil)
			},
			expectError: true,
		},
		{
			name: "Should return err on service err",
			rec:  httptest.NewRecorder(),
			reqFn: func() *http.Request {
				return httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			},
			mockFn: func() *mock.Call {
				return serviceMock.
					On("Create", mock.Anything).
					Return(errors.New("fail"))
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockCall *mock.Call
			if tt.mockFn != nil {
				mockCall = tt.mockFn()
			}

			fn := handler.CreateUser()
			fn(tt.rec, tt.reqFn())

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, tt.rec.Code)
				return
			}

			if tt.expectResult {
				assert.Equal(t, http.StatusOK, tt.rec.Code)
			}
			serviceMock.AssertExpectations(t)
			mockCall.Unset()
		})
	}
}

func Test_userHandler_Login(t *testing.T) {
	serviceMock := new(mocks.UserServiceMock)
	handler := userHandler{
		service: serviceMock,
	}

	tests := []struct {
		name         string
		rec          *httptest.ResponseRecorder
		reqFn        func() *http.Request
		mockFn       func() *mock.Call
		expectResult bool
		expectError  bool
	}{
		//{
		//	name: "Should login",
		//	rec:  httptest.NewRecorder(),
		//	reqFn: func() *http.Request {
		//		return httptest.NewRequest("POST", "/", strings.NewReader("{}"))
		//	},
		//	mockFn: func() *mock.Call {
		//		return serviceMock.
		//			On("GetByExternalId", mock.Anything).
		//			Return(user, nil)
		//	},
		//	expectResult: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockCall *mock.Call
			if tt.mockFn != nil {
				mockCall = tt.mockFn()
			}

			fn := handler.Login()
			fn(tt.rec, tt.reqFn())

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, tt.rec.Code)
				return
			}

			if tt.expectResult {
				assert.Equal(t, http.StatusOK, tt.rec.Code)
				assert.NotNil(t, tt.rec.Body)
			}
			serviceMock.AssertExpectations(t)
			mockCall.Unset()
		})
	}
}
