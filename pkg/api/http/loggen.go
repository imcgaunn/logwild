package http

import (
	"io"
	"log/slog"
	"net/http"
	"os"

	"mcgaunn.com/logwild/pkg/logmaker"
)

// Loggen godoc
// @Summary Log generation endpoint
// @Description starts logging messages and reports stats
// @Tags HTTP API
// @Accept json
// @Produce json
// @Success 200 {object} api.LogStatsResponse
// @Router /api/loggen [get]
func (s *Server) logGenHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "logGenHandler")
	defer span.End()
	lvl := &slog.LevelVar{}
	fp, err := os.Create("/tmp/loggen.log")
	stdoutAndFile := io.MultiWriter(fp, os.Stdout)
	if err != nil {
		s.logger.Error("failed to create log file: ", err)
		panic(err)
	}
	h := slog.NewJSONHandler(stdoutAndFile, &slog.HandlerOptions{Level: lvl})
	lm := logmaker.NewLogMaker(logmaker.WithLogger(slog.New(h)))
	donech := make(chan int)
	go func() {
		err := lm.StartWriting(donech)
		if err != nil {
			panic(err) // something bad here
		}
	}()
	logCount := <-donech
	data := LogStatsResponse{logCount: logCount}
	s.JSONResponse(w, r, data)
}

type LogStatsResponse struct {
	logCount int
}
