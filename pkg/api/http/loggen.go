package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

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
	queryVals := r.URL.Query()
	_, span := s.tracer.Start(r.Context(), "logGenHandler")
	defer span.End()
	lvl := &slog.LevelVar{}
	perSecondRateParam := queryVals.Get("per_second")
	perMessageSizeParam := queryVals.Get("message_size")
	burstDuration := queryVals.Get("burst_dur")

	var fp *os.File
	if s.config.LogwildOutFile == "-" {
		fp = os.Stdout
	} else {
		var err error
		fp, err = os.Create(s.config.LogwildOutFile)
		if err != nil {
			s.logger.Error("failed to create log file: ", err, "fileName", s.config.LogwildOutFile)
			panic(err)
		}
	}
	h := slog.NewJSONHandler(fp, &slog.HandlerOptions{Level: lvl})
	// depepending on values passed in through query string, construct
	// a new log maker with the appropriate options.
	var optFuncs []logmaker.OptFunc
	if perSecondRateParam != "" {
		fmt.Printf("got persecondrate %s\n", perSecondRateParam)
		perSecondRateInt, err := strconv.ParseInt(perSecondRateParam, 10, 64)
		if err != nil {
			panic(err) // could not convert perSecondRate param to integer :(
		}
		optFuncs = append(optFuncs, logmaker.WithPerSecondRate(perSecondRateInt))
	}
	if perMessageSizeParam != "" {
		fmt.Printf("got permessagesize %s\n", perMessageSizeParam)
		perMessageSizeInt, err := strconv.ParseInt(perMessageSizeParam, 10, 64)
		if err != nil {
			panic(err) // could not convert perMessageSize param to integer :(
		}
		optFuncs = append(optFuncs, logmaker.WithPerMessageSizeBytes(perMessageSizeInt))
	}
	if burstDuration != "" {
		fmt.Printf("got burstduration %s\n", burstDuration)
		burstDurationInt, err := strconv.ParseInt(burstDuration, 10, 0)
		if err != nil {
			panic(err)
		}
		optFuncs = append(optFuncs, logmaker.WithBurstDuration(time.Second*time.Duration(burstDurationInt)))
	}
	optFuncs = append(optFuncs, logmaker.WithLogger(slog.New(h)))
	lm := logmaker.NewLogMaker(optFuncs...)
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
