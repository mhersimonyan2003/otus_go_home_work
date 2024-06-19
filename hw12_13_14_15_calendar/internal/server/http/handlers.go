package internalhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Failed to decode request: " + err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event := storage.Event{
		ID:        req.ID,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Details:   req.Details,
	}

	if err := s.app.AddEvent(event); err != nil {
		s.logger.Error("Failed to add event: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Failed to decode request: " + err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event := storage.Event{
		ID:        req.ID,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Details:   req.Details,
	}

	if err := s.app.UpdateEvent(event); err != nil {
		s.logger.Error("Failed to update event: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]

	if err := s.app.DeleteEvent(eventID); err != nil {
		s.logger.Error("Failed to delete event: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	start, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	if err != nil {
		http.Error(w, "Invalid end time", http.StatusBadRequest)
		return
	}

	events, err := s.app.ListEvents(start, end)
	if err != nil {
		s.logger.Error("Failed to list events: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(events); err != nil {
		s.logger.Error("Failed to encode response: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]

	event, err := s.app.GetEventByID(eventID)
	if err != nil {
		s.logger.Error("Failed to get event: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(event); err != nil {
		s.logger.Error("Failed to encode response: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) helloWorldHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, World!"))
}
