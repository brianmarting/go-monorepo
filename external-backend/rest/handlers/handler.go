package handlers

import (
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
	h.Use(otelchi.Middleware("external-backend", otelchi.WithChiRoutes(h)))

	h.Route("/coal", func(router chi.Router) {
		router.Post("/", h.mineralHandler.Post())
	})
}
