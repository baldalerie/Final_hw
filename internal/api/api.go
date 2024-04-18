package api

import (
	"net"
	"net/http"

	"FinalTaskAppGoBasic/internal/configs"
	"FinalTaskAppGoBasic/internal/handlers"
	"FinalTaskAppGoBasic/internal/logs"

	"github.com/gorilla/mux"
)

type Server struct {
	address string

	router   *mux.Router
	handlers *handlers.Handlers
}

func (s *Server) Init() {
	// s.router.HandleFunc("/transactions", handlers.HandleTransactions).Methods("GET", "POST")
	s.router.HandleFunc("/transactions/{id}", s.handlers.HandleTransactions).Methods("GET")
	s.router.HandleFunc("/transactions/{id}", s.handlers.HandleTransactions).Methods("PUT, DELETE")

	s.router.HandleFunc("/users", s.handlers.RegisterUser).Methods("POST")
	s.router.HandleFunc("/users/login", s.handlers.LoginUser).Methods("POST")
}

func (s *Server) ListenAndServe() error {
	logs.Log.
		WithField("address", s.address).
		Info("server start listening")

	return http.ListenAndServe(s.address, s.router)
}

func New(cfg *configs.Api, apiHandlers *handlers.Handlers) *Server {
	return &Server{
		address:  net.JoinHostPort(cfg.Host, cfg.Port),
		router:   mux.NewRouter(),
		handlers: apiHandlers,
	}
}
