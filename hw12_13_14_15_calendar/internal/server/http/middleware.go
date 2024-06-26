package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(logg logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			latency := time.Since(start)

			logg.Info(fmt.Sprintf(
				"%s [%s] %s %s %s %d %s %q",
				r.RemoteAddr,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.RequestURI,
				r.Proto,
				r.Response.StatusCode,
				latency,
				r.UserAgent(),
			))
		})
	}
}
