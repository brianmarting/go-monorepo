package mocks

import (
	"common/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MineralPublisherMock struct {
	mock.Mock
}

func (m *MineralPublisherMock) Publish(ctx context.Context, model model.Mineral) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}
