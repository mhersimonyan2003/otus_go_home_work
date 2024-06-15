package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_AddEvent(t *testing.T) {
	store := New()
	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Details:   "This is a test event",
	}

	err := store.AddEvent(event)
	require.NoError(t, err, "unexpected error adding event")
}

func TestMemoryStorage_AddEvent_DateBusy(t *testing.T) {
	store := New()
	startTime := time.Now()
	event1 := storage.Event{
		ID:        "1",
		Title:     "Test Event 1",
		StartTime: startTime,
		EndTime:   startTime.Add(time.Hour),
		Details:   "This is test event 1",
	}
	event2 := storage.Event{
		ID:        "2",
		Title:     "Test Event 2",
		StartTime: startTime.Add(time.Minute * 30),
		EndTime:   startTime.Add(time.Hour * 2),
		Details:   "This is test event 2",
	}

	err := store.AddEvent(event1)
	require.NoError(t, err, "unexpected error adding event 1")

	err = store.AddEvent(event2)
	require.ErrorIs(t, err, storage.ErrDateBusy, "expected ErrDateBusy")
}

func TestMemoryStorage_UpdateEvent(t *testing.T) {
	store := New()
	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Details:   "This is a test event",
	}

	err := store.AddEvent(event)
	require.NoError(t, err, "unexpected error adding event")

	event.Title = "Updated Test Event"
	err = store.UpdateEvent(event)
	require.NoError(t, err, "unexpected error updating event")

	updatedEvent, err := store.GetEventByID(event.ID)
	require.NoError(t, err, "unexpected error getting event by ID")
	require.Equal(t, "Updated Test Event", updatedEvent.Title, "event title was not updated")
}

func TestMemoryStorage_UpdateEvent_NotFound(t *testing.T) {
	store := New()
	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Details:   "This is a test event",
	}

	err := store.UpdateEvent(event)
	require.ErrorIs(t, err, storage.ErrEventNotFound, "expected ErrEventNotFound")
}

func TestMemoryStorage_DeleteEvent(t *testing.T) {
	store := New()
	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Details:   "This is a test event",
	}

	err := store.AddEvent(event)
	require.NoError(t, err, "unexpected error adding event")

	err = store.DeleteEvent(event.ID)
	require.NoError(t, err, "unexpected error deleting event")

	_, err = store.GetEventByID(event.ID)
	require.ErrorIs(t, err, storage.ErrEventNotFound, "expected ErrEventNotFound")
}

func TestMemoryStorage_DeleteEvent_NotFound(t *testing.T) {
	store := New()
	err := store.DeleteEvent("non-existent-id")
	require.ErrorIs(t, err, storage.ErrEventNotFound, "expected ErrEventNotFound")
}

func TestMemoryStorage_ListEvents(t *testing.T) {
	store := New()
	startTime := time.Now()

	event1 := storage.Event{
		ID:        "1",
		Title:     "Test Event 1",
		StartTime: startTime,
		EndTime:   startTime.Add(time.Hour),
		Details:   "This is test event 1",
	}
	event2 := storage.Event{
		ID:        "2",
		Title:     "Test Event 2",
		StartTime: startTime.Add(time.Hour),
		EndTime:   startTime.Add(time.Hour * 2),
		Details:   "This is test event 2",
	}

	err := store.AddEvent(event1)
	require.NoError(t, err, "unexpected error adding event 1")
	err = store.AddEvent(event2)
	require.NoError(t, err, "unexpected error adding event 2")

	events, err := store.ListEvents(startTime, startTime.Add(time.Hour*2))
	require.NoError(t, err, "unexpected error listing events")
	require.Len(t, events, 2, "expected 2 events")
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	store := New()
	startTime := time.Now()

	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: startTime,
		EndTime:   startTime.Add(time.Hour),
		Details:   "This is a test event",
	}

	err := store.AddEvent(event)
	require.NoError(t, err, "unexpected error adding event")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = store.ListEvents(startTime, startTime.Add(time.Hour*2))
		}()
	}

	wg.Wait()
}
