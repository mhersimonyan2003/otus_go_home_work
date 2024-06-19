package app

import (
	"time"

	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.Storage
}

func New(storage storage.Storage) *App {
	return &App{storage: storage}
}

func (a *App) AddEvent(event storage.Event) error {
	return a.storage.AddEvent(event)
}

func (a *App) UpdateEvent(event storage.Event) error {
	return a.storage.UpdateEvent(event)
}

func (a *App) DeleteEvent(eventID string) error {
	return a.storage.DeleteEvent(eventID)
}

func (a *App) ListEvents(start, end time.Time) ([]storage.Event, error) {
	return a.storage.ListEvents(start, end)
}

func (a *App) GetEventByID(eventID string) (storage.Event, error) {
	return a.storage.GetEventByID(eventID)
}
