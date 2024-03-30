package httpx

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)

		next.ServeHTTP(w, r)
	})
}
