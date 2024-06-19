package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	logger := logger.New(logger.Debug)
	storage := memorystorage.New()
	app := app.New(storage)
	server := NewServer(logger, app, "localhost", 8080)

	tests := []struct {
		name       string
		method     string
		target     string
		body       string
		handler    func(http.ResponseWriter, *http.Request)
		vars       map[string]string
		params     map[string]string
		wantStatus int
		wantBody   string
	}{
		{
			name:   "CreateEvent",
			method: http.MethodPost,
			target: "/events",
			body: `{
				"id":         "1",
				"title":      "Created Test Event",
				"start_time": "2023-06-19T10:00:00Z",
				"end_time":   "2023-06-19T11:00:00Z",
				"details":    "This is initital details of the test event"
			}`,
			handler:    server.createEventHandler,
			wantStatus: http.StatusCreated,
			wantBody:   "",
		},
		{
			name:   "CreateEvent",
			method: http.MethodPost,
			target: "/events",
			body: `{
				"id":         "1",
				"title":      "Created Test Event",
				"start_time": "2023-06-19T10:00:00Z",
				"end_time":   "2023-06-19T11:00:00Z",
				"details":    "This is initital details of the test event"
			}`,
			handler:    server.createEventHandler,
			wantStatus: http.StatusInternalServerError,
			wantBody:   "the time slot is already occupied by another event\n",
		},
		{
			name:   "UpdateEvent",
			method: http.MethodPut,
			target: "/events/1",
			body: `{
				"id":         "1",
				"title":      "Updated Test Event",
				"start_time": "2023-06-19T10:00:00Z",
				"end_time":   "2023-06-19T11:00:00Z",
				"details":    "This is an updated test event details"
			}`,
			handler:    server.updateEventHandler,
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "ListEvents",
			method:     http.MethodGet,
			target:     "/events",
			body:       "",
			handler:    server.listEventsHandler,
			params:     map[string]string{"start": "2023-06-19T00:00:00Z", "end": "2023-06-20T00:00:00Z"},
			wantStatus: http.StatusOK,
			wantBody: "[{\"ID\":\"1\",\"Title\":\"Updated Test Event\",\"StartTime\":\"2023-06-19T10:00:00Z\"," +
				"\"EndTime\":\"2023-06-19T11:00:00Z\",\"Details\":\"This is an updated test event details\"}]\n",
		},
		{
			name:       "GetEvent",
			method:     http.MethodGet,
			target:     "/events/1",
			body:       "",
			handler:    server.getEventHandler,
			vars:       map[string]string{"id": "1"},
			wantStatus: http.StatusOK,
			wantBody: "{\"ID\":\"1\",\"Title\":\"Updated Test Event\",\"StartTime\":\"2023-06-19T10:00:00Z\"," +
				"\"EndTime\":\"2023-06-19T11:00:00Z\",\"Details\":\"This is an updated test event details\"}\n",
		},
		{
			name:       "DeleteEvent",
			method:     http.MethodDelete,
			target:     "/events/1",
			body:       "",
			handler:    server.deleteEventHandler,
			vars:       map[string]string{"id": "1"},
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "ListEvents",
			method:     http.MethodGet,
			target:     "/events",
			body:       "",
			handler:    server.listEventsHandler,
			params:     map[string]string{"start": "2023-06-19T00:00:00Z", "end": "2023-06-20T00:00:00Z"},
			wantStatus: http.StatusOK,
			wantBody:   "null\n",
		},
		{
			name:       "GetEvent",
			method:     http.MethodGet,
			target:     "/events/1",
			body:       "",
			handler:    server.getEventHandler,
			vars:       map[string]string{"id": "1"},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "event not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mux.NewRouter()

			r.HandleFunc(tt.target, func(w http.ResponseWriter, r *http.Request) {
				if tt.vars != nil {
					r = mux.SetURLVars(r, tt.vars)
				}
				tt.handler(w, r)
			}).Methods(tt.method)

			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))

			q := req.URL.Query()
			for k, v := range tt.params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantBody, rr.Body.String())
		})
	}
}

func TestStartAndStop(t *testing.T) {
	logger := logger.New(logger.Debug)
	storage := memorystorage.New()
	app := app.New(storage)
	server := NewServer(logger, app, "localhost", 8080)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)

	go func() {
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	time.Sleep(1 * time.Second)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server failed to start: %v", err)
		}
	default:
	}

	if err := server.Stop(context.Background()); err != nil {
		t.Fatalf("server failed to stop: %v", err)
	}
}
