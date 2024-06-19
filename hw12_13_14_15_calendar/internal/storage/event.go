package storage

import (
	"errors"
	"time"
)

// Event структура события.
type Event struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Details   string    `db:"details"`
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
