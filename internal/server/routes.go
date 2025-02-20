package server

import (
	"github.com/bismastr/scrapper-example/internal/mrt"
	"github.com/gorilla/mux"
)

func (s *Server) SetupRoutes(handler *mrt.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/schedules", handler.GetAllStation).Methods("GET")
	router.HandleFunc("/schedulesById", handler.GetScheduleById).Methods("GET")
	router.HandleFunc("/stations_name", handler.GetAllStationName).Methods("GET")
	return router
}
