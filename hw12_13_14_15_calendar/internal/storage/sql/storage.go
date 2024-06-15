package sqlstorage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

// SQLStorage реализация интерфейса Storage для PostgreSQL.
type SQLStorage struct {
	db *sqlx.DB
}

// New создает новый экземпляр SQLStorage.
func New(dsn string) (*SQLStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLStorage{db: db}, nil
}

func (s *SQLStorage) AddEvent(event storage.Event) error {
	query := `INSERT INTO events (id, title, start_time, end_time, details)
		VALUES (:id, :title, :start_time, :end_time, :details)`
	_, err := s.db.NamedExec(query, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) UpdateEvent(event storage.Event) error {
	query := `UPDATE events SET title=:title, start_time=:start_time, end_time=:end_time, details=:details WHERE id=:id`
	_, err := s.db.NamedExec(query, event)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) DeleteEvent(eventID string) error {
	query := `DELETE FROM events WHERE id=$1`
	_, err := s.db.Exec(query, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) ListEvents(start, end time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, start_time, end_time, details FROM events WHERE start_time < $2 AND end_time > $1`
	var events []storage.Event
	err := s.db.Select(&events, query, start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) GetEventByID(eventID string) (storage.Event, error) {
	query := `SELECT id, title, start_time, end_time, details FROM events WHERE id=$1`
	var event storage.Event
	err := s.db.Get(&event, query, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrEventNotFound
		}
		return storage.Event{}, err
	}
	return event, nil
}
