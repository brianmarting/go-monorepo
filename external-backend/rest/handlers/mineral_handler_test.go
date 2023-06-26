package handlers

import (
	"common/model"
	"errors"
	"external-backend/service/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mineral = model.Mineral{}

func Test_mineralHandler_Post(t *testing.T) {
	serviceMock := new(mocks.MineralServiceMock)
	handler := mineralHandler{service: serviceMock}

	tests := []struct {
		name         string
		rec          *httptest.ResponseRecorder
		reqFn        func() *http.Request
		mockFn       func() *mock.Call
		expectResult bool
		expectError  bool
	}{
		{
			name: "Should add mineral",
			rec:  httptest.NewRecorder(),
			reqFn: func() *http.Request {
				return httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			},
			mockFn: func() *mock.Call {
				return serviceMock.On("AddMineral", mock.Anything, mineral).Return(nil)
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
			name: "Should return err on add mineral err",
			rec:  httptest.NewRecorder(),
			reqFn: func() *http.Request {
				return httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			},
			mockFn: func() *mock.Call {
				return serviceMock.On("AddMineral", mock.Anything, mineral).Return(errors.New("failed to add mineral"))
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
			defer func() {
				if mockCall != nil {
					mockCall.Unset()
				}
			}()

			fn := handler.Post()
			fn(tt.rec, tt.reqFn())

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, tt.rec.Code)
				return
			}

			assert.Equal(t, http.StatusOK, tt.rec.Code)
			serviceMock.AssertExpectations(t)
		})
	}
}
