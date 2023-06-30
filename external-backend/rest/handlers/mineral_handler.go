package handlers

import (
	"encoding/json"
	"go-monorepo/external-backend/service"
	"go-monorepo/internal/model"
	"net/http"
)

type MineralHandler interface {
	Post() http.HandlerFunc
}

type mineralHandler struct {
	service service.MineralService
}

func NewMineralHandler() MineralHandler {
	return &mineralHandler{
		service: service.GetMineralService(),
	}
}

func (h mineralHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mineral model.Mineral

		if err := json.NewDecoder(r.Body).Decode(&mineral); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.service.AddMineral(r.Context(), mineral); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(200)
	}
}
