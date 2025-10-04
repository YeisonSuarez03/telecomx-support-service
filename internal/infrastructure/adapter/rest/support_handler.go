package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"telecomx-support-service/internal/application/service"
	"telecomx-support-service/internal/domain/model"
)

type SupportHandler struct {
	service *service.SupportService
}

func NewSupportHandler(s *service.SupportService) *SupportHandler {
	return &SupportHandler{service: s}
}

func (h *SupportHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/support", h.handleSupport)
}

func (h *SupportHandler) handleSupport(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		data, err := h.service.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)

	case http.MethodPost:
		var p model.Support
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.service.Create(ctx, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
