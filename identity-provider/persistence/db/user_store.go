package db

import (
	"go-monorepo/internal/model"
)

type UserStore interface {
	GetByExternalId(name string) (model.User, error)
	Create(user model.User) error
}
