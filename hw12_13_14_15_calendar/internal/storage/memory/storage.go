package memorystorage

import (
	"sync"
	"time"

	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

// MemoryStorage реализация интерфейса Storage на памяти.
type MemoryStorage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

// New создает новый экземпляр MemoryStorage.
func New() *MemoryStorage {
	return &MemoryStorage{
		events: make(map[string]storage.Event),
	}
}

func (s *MemoryStorage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if event.StartTime.Before(e.EndTime) && event.EndTime.After(e.StartTime) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *MemoryStorage) UpdateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventNotFound
	}

	s.events[event.ID] = event
	return nil
}

func (s *MemoryStorage) DeleteEvent(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return storage.ErrEventNotFound
	}

	delete(s.events, eventID)
	return nil
}

func (s *MemoryStorage) ListEvents(start, end time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemoryStorage) GetEventByID(eventID string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[eventID]
	if !ok {
		return storage.Event{}, storage.ErrEventNotFound
	}

	return event, nil
}
