package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	app        Application
}

type Application interface {
	// Define methods that will be used by the application
}

func NewServer(logger logger.Logger, app Application) *Server {
	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger))
	r.HandleFunc("/", helloWorldHandler).Methods(http.MethodGet)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting server on " + s.httpServer.Addr)
	errCh := make(chan error)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		s.logger.Info("Server exited properly")
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
