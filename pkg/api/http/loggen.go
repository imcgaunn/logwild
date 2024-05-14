package http

import "net/http"

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

	data := LogStatsResponse{}
	s.JSONResponse(w, r, data)
}

type LogStatsResponse struct{}
