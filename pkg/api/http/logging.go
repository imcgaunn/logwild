package http

import (
	"log/slog"
	"net/http"
)

type LoggingMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug(
			"request started",
			slog.String("proto", r.Proto),
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.String("remote", r.RemoteAddr),
			slog.String("user-agent", r.UserAgent()),
		)
		next.ServeHTTP(w, r)
	})
}
