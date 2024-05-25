package http

import (
	"log/slog"
	"os"
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
		LogwildOutFile:        "-",
		LogwildPerSecondRate:  5000,
		LogwildPerMessageSize: 50,
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))
	logger := slog.Default().With("mockserver", "yes")
	return &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
		tracer: noop.NewTracerProvider().Tracer("mock"),
	}
}
