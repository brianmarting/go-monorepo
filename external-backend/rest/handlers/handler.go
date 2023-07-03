package handlers

import (
	"go-monorepo/external-backend/rest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	CreateAllRoutes()
}

type handler struct {
	*chi.Mux

	mineralHandler MineralHandler
}

func NewHandler() Handler {
	return &handler{
		Mux:            chi.NewMux(),
		mineralHandler: NewMineralHandler(),
	}
}

func (h *handler) CreateAllRoutes() {
	h.Use(otelchi.Middleware("go-monorepo/external-backend", otelchi.WithChiRoutes(h)))
	h.Use(middleware.IsAuthenticated())

	h.Route("/mineral", func(router chi.Router) {
		router.Post("/deposit", h.mineralHandler.Post())
	})
}
