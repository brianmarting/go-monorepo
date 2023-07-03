package service

import (
	"errors"
	"go-monorepo/identity-provider/persistence/db/mocks"
	"go-monorepo/internal/model"
	"testing"

	"github.com/stretchr/testify/mock"
)

var user = model.User{
	Id:           32,
	ExternalId:   "dq9b5f7e-61e6-480e-91bb-269f7a929622",
	Name:         "John",
	Password:     "somepwd",
	TokenVersion: 1,
}

func Test_userService_GetByExternalId(t *testing.T) {
	storeMock := new(mocks.UserStoreMock)
	service := userService{
		store: storeMock,
	}

	tests := []struct {
		name    string
		arg     string
		mockFn  func() *mock.Call
		wantErr bool
	}{
		{
			name: "Should get user by external id",
			arg:  user.ExternalId,
			mockFn: func() *mock.Call {
				return storeMock.
					On("GetByExternalId", user.ExternalId).
					Return(user, nil)
			},
		},
		{
			name: "Should return err",
			arg:  user.ExternalId,
			mockFn: func() *mock.Call {
				return storeMock.
					On("GetByExternalId", user.ExternalId).
					Return(model.User{}, errors.New("fail"))
			},
			wantErr: true,
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

			result, err := service.GetByExternalId(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByExternalId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && result != user {
				t.Errorf("Bodies dont match")
			}
		})
	}
}

func Test_userService_Create(t *testing.T) {
	storeMock := new(mocks.UserStoreMock)
	service := userService{
		store: storeMock,
	}

	tests := []struct {
		name    string
		arg     model.User
		mockFn  func() *mock.Call
		wantErr bool
	}{
		{
			name: "Should create user",
			arg:  user,
			mockFn: func() *mock.Call {
				return storeMock.
					On("Create", user).
					Return(nil)
			},
		},
		{
			name: "Should return err",
			arg:  user,
			mockFn: func() *mock.Call {
				return storeMock.
					On("Create", user).
					Return(errors.New("fail"))
			},
			wantErr: true,
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

			err := service.Create(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
