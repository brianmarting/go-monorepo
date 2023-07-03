package service

import (
	"go-monorepo/identity-provider/persistence/db"
	"go-monorepo/identity-provider/persistence/db/psql"
	"go-monorepo/internal/model"
)

type UserService interface {
	GetByExternalId(id string) (model.User, error)
	Create(user model.User) error
}

type userService struct {
	store db.UserStore
}

func NewUserService() UserService {
	return &userService{
		store: psql.NewUserStore(),
	}
}

func (s userService) GetByExternalId(id string) (model.User, error) {
	return s.store.GetByExternalId(id)
}

func (s userService) Create(user model.User) error {
	return s.store.Create(user)
}
