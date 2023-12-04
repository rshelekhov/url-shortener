package logger

import (
	"github.com/go-chi/chi/middleware"
	mw "github.com/rshelekhov/url-shortener/internal/http-server/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		// Handler
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Get information about the request
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String(mw.RequestID, middleware.GetReqID(r.Context())),
			)

			// Create ResponseWriter for getting details about the response
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Time when we get a response
			t1 := time.Now()

			// Add details to the log in defer
			// At this point the request will already be processed
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Pass control to the next handler in the middleware chain
			next.ServeHTTP(ww, r)
		}

		// Return the handler created above by casting it to type http.HandlerFunc
		return http.HandlerFunc(fn)
	}
}
