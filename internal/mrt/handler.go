package mrt

import (
	"fmt"
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
	fmt.Println("Test GetAllStation")
	h.mrtService.GetAllStation(r.Context())
}
