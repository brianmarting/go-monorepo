package mocks

import (
	"go-monorepo/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetByExternalId(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserServiceMock) Create(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
