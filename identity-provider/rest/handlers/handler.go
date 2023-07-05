package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	CreateAllRoutes()
}

type handler struct {
	*chi.Mux

	userHandler  UserHandler
	tokenHandler TokenHandler
}

func NewHandler() Handler {
	return &handler{
		Mux:          chi.NewMux(),
		userHandler:  NewUserHandler(),
		tokenHandler: NewTokenHandler(),
	}
}

func (h handler) CreateAllRoutes() {
	h.Route("/user", func(router chi.Router) {
		router.Post("/login", h.userHandler.Login())
		router.Post("/create", h.userHandler.CreateUser())
		router.Route("/token", func(router chi.Router) {
			router.Post("/validate", h.tokenHandler.Validate())
			router.Post("/refresh", h.tokenHandler.PostRefreshToken())
		})
	})
}
