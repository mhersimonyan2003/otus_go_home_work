package internalhttp

import (
	"time"
)

type EventRequest struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Details   string    `json:"details"`
}

type EventResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Details   string    `json:"details"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
