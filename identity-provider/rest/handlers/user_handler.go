package handlers

import (
	"encoding/json"
	"errors"
	"go-monorepo/identity-provider/service"
	"go-monorepo/internal/model"
	"net/http"

	"github.com/rs/zerolog/log"
)

type UserHandler interface {
	CreateUser() http.HandlerFunc
	Login() http.HandlerFunc
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler() UserHandler {
	return &userHandler{
		service: service.NewUserService(),
	}
}

func (h userHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.service.Create(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(200)
	}
}

func ResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

func (h userHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := h.service.GetByExternalId(user.ExternalId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if dbUser.Password != user.Password {
			http.Error(w, errors.New("invalid password").Error(), http.StatusBadRequest)
			return
		}

		accessToken, err := CreateAccessToken(dbUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error().Err(err).Msg("failed to sign token")
			return
		}

		http.SetCookie(w, CreateCookie(dbUser))
		if _, err = w.Write(accessToken); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
