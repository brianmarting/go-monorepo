package db

import (
	"go-monorepo/internal/model"
)

type MineralStore interface {
	GetByName(name string) (model.Mineral, error)
	UpdateAmount(name string, amount int) error
}
