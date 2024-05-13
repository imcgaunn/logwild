package signals

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Shutdown struct {
	logger                *slog.Logger
	tracerProvider        *sdktrace.TracerProvider
	serverShutdownTimeout time.Duration
}

func NewShutdown(serverShutdownTimeout time.Duration, logger *slog.Logger) (*Shutdown, error) {
	srv := &Shutdown{
		logger:                logger,
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv, nil
}

func (s *Shutdown) Graceful(stopCh <-chan struct{}, httpServer *http.Server, healthy *int32, ready *int32) {
	ctx := context.Background()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	s.logger.Debug("got signal, shutting down")

	// all calls to /healthz and /readyz will fail from now on
	atomic.StoreInt32(healthy, 0)
	atomic.StoreInt32(ready, 0)

	s.logger.Debug("set 'healthy' and 'ready' indicators 0", "healthy", healthy, "ready", ready)
	s.logger.Info("Shutting down HTTP/HTTPS server", slog.Duration("timeout", s.serverShutdownTimeout))

	// There could be a period where a terminating pod may still receive requests. Implementing a brief wait can mitigate this.
	// See: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination
	// the readiness check interval must be lower than the timeout
	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}
	// stop OpenTelemetry tracer provider
	if s.tracerProvider != nil {
		if err := s.tracerProvider.Shutdown(ctx); err != nil {
			s.logger.Warn("stopping tracer provider", "err", err)
		}
	}
	// determine if the http server was started
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			s.logger.Warn("HTTP server graceful shutdown failed", "err", err)
		}
	}
}
