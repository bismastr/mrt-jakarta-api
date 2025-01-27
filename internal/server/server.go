package server

import (
	"net/http"
	"time"
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

func (s *Server) Start() {
	s.ListenAndServe()
}
