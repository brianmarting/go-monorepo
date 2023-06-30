package mocks

import (
	"context"
	"go-monorepo/internal/model"

	"github.com/stretchr/testify/mock"
)

type MineralServiceMock struct {
	mock.Mock
}

func (m *MineralServiceMock) AddMineral(ctx context.Context, mineral model.Mineral) error {
	args := m.Called(ctx, mineral)
	return args.Error(0)
}
