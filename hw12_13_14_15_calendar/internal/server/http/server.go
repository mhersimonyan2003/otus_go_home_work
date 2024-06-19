package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	app        *app.App
}

func NewServer(logger logger.Logger, app *app.App, host string, port int) *Server {
	s := &Server{
		logger: logger,
		app:    app,
	}

	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger))
	r.HandleFunc("/", s.helloWorldHandler).Methods(http.MethodGet)
	r.HandleFunc("/events", s.createEventHandler).Methods(http.MethodPost)
	r.HandleFunc("/events/{id}", s.updateEventHandler).Methods(http.MethodPut)
	r.HandleFunc("/events/{id}", s.deleteEventHandler).Methods(http.MethodDelete)
	r.HandleFunc("/events", s.listEventsHandler).Methods(http.MethodGet)
	r.HandleFunc("/events/{id}", s.getEventHandler).Methods(http.MethodGet)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP server on " + s.httpServer.Addr)
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
