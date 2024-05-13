package http

import (
	"log/slog"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace/noop"
)

func NewMockServer() *Server {
	config := &Config{
		Port:                  "9998",
		ServerShutdownTimeout: 5 * time.Second,
		HttpServerTimeout:     30 * time.Second,
		BackendURL:            []string{},
		DataPath:              "/data",
		ConfigPath:            "/config",
		HttpClientTimeout:     30 * time.Second,
		Hostname:              "localhost",
	}

	logger := slog.Default().With()
	return &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
		tracer: noop.NewTracerProvider().Tracer("mock"),
	}
}
