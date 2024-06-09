package storage

import (
	"errors"
	"time"
)

// Event структура события.
type Event struct {
	ID        string
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Details   string
}

// Storage интерфейс для работы с хранилищем событий.
type Storage interface {
	AddEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(eventID string) error
	ListEvents(start, end time.Time) ([]Event, error)
	GetEventByID(eventID string) (Event, error)
}

// Ошибки, связанные с бизнес-логикой.
var (
	ErrDateBusy      = errors.New("the time slot is already occupied by another event")
	ErrEventNotFound = errors.New("event not found")
)
