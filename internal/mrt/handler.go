package mrt

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (h *Handler) GetScheduleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Printf("Missing id")
	}
	isHoliday, err := strconv.ParseBool(r.URL.Query().Get("is_holiday"))
	if err != nil {
		log.Printf("Missing id")
	}
	directionStationId, err := strconv.ParseInt(r.URL.Query().Get("direction_station_id"), 10, 64)
	if err != nil {
		log.Printf("Missing id")
	}

	result := h.mrtService.GetScheduleById(r.Context(), id, isHoliday, directionStationId)

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAllStationName(w http.ResponseWriter, r *http.Request) {
	stations := h.mrtService.GetAllStationName(r.Context())

	json.NewEncoder(w).Encode(stations)
}
