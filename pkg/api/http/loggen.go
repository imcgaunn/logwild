package http

import (
	"log/slog"
	"mcgaunn.com/logwild/pkg/logmaker"
	"net/http"
	"os"
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
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
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
