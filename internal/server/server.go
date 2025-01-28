package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server
}

func NewServer() *Server {
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		Server: s,
	}
}

func (s *Server) Start(router *mux.Router) {
	s.Handler = router
	fmt.Println("Server started on port :8080")
	s.ListenAndServe()
}
