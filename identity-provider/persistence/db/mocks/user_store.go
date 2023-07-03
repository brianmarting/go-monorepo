package mocks

import (
	"go-monorepo/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) GetByExternalId(name string) (model.User, error) {
	args := m.Called(name)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserStoreMock) Create(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
