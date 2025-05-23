package http

import (
	"net/http"
	"runtime"
	"strconv"

	"mcgaunn.com/logwild/pkg/version"
)

// Info godoc
// @Summary Runtime information
// @Description returns the runtime information
// @Tags HTTP API
// @Accept json
// @Produce json
// @Success 200 {object} api.RuntimeResponse
// @Router /api/info [get]
func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "infoHandler")
	defer span.End()

	data := RuntimeResponse{
		Hostname:     s.config.Hostname,
		Version:      version.VERSION,
		Revision:     version.REVISION,
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		Runtime:      runtime.Version(),
		NumGoroutine: strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		NumCPU:       strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}

	s.JSONResponse(w, r, data)
}

type RuntimeResponse struct {
	Hostname     string `json:"hostname"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	Runtime      string `json:"runtime"`
	NumGoroutine string `json:"num_goroutine"`
	NumCPU       string `json:"num_cpu"`
}
