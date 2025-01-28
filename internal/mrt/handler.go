package mrt

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	mrtService *MrtService
}

func NewHandler(mrt *MrtService) *Handler {
	return &Handler{
		mrtService: mrt,
	}
}

func (h *Handler) GetAllStation(w http.ResponseWriter, r *http.Request) {

	result := h.mrtService.GetAllStation(r.Context())
	json.NewEncoder(w).Encode(result)

}
