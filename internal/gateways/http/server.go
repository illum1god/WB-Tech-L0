package http

import (
	"WB-L0/internal/usecase"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type Server struct {
	host       string
	port       uint16
	router     *gin.Engine
	httpServer *http.Server
	logger     *logrus.Logger
}

func NewServer(service usecase.Service, options ...func(*Server)) *Server {
	r := gin.Default()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	s := &Server{
		router: r,
		host:   "localhost",
		port:   8080,
		logger: logger,
	}

	for _, o := range options {
		o(s)
	}

	setupRouter(r, service, logger)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.router,
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server ListenAndServe: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	log.Println("Server gracefully stopped")
	return nil
}
