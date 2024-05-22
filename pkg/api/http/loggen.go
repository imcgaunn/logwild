package http

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	_, span := s.tracer.Start(r.Context(), "logGenHandler")
	defer span.End()
	span.AddEvent("startInitializeLogger")
	h := s.createLogHandlerOrPanic()
	// create initial options from config
	optFuncs := s.buildLoggerOptionsFromConfig()
	// override functions based on query params
	optFuncs = append(optFuncs, s.buildLoggerOptionsFromQueryParams(h, r)...)
	lm := logmaker.NewLogMaker(optFuncs...)
	span.AddEvent("doneInitializeLogger")
	slog.Info("lm config", "perSecondRate", lm.PerSecondRate)
	donech := make(chan int)
	go func() {
		err := lm.StartWriting(donech)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			panic(err) // something bad here
		}
	}()
	span.AddEvent("startedWriting", trace.WithAttributes(attribute.Int("logCount", 0)))
	logCount := <-donech
	data := LogStatsResponse{logCount: logCount}
	span.AddEvent("doneWriting", trace.WithAttributes(attribute.Int("logCount", logCount)))
	span.SetStatus(codes.Ok, "successfully wrote logs")
	s.JSONResponse(w, r, data)
}

func (s *Server) createLogHandlerOrPanic() *slog.JSONHandler {
	lvl := &slog.LevelVar{}

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
	return h
}

func (s *Server) buildLoggerOptionsFromConfig() []logmaker.OptFunc {
	var optFuncs []logmaker.OptFunc
	h := s.createLogHandlerOrPanic()
	optFuncs = append(optFuncs, logmaker.WithLogger(slog.New(h)))
	optFuncs = append(optFuncs, logmaker.WithPerSecondRate(s.config.LogwildPerSecondRate))
	return optFuncs
}

func (s *Server) buildLoggerOptionsFromQueryParams(h *slog.JSONHandler, r *http.Request) []logmaker.OptFunc {
	var optFuncs []logmaker.OptFunc
	queryVals := r.URL.Query()
	_, span := s.tracer.Start(r.Context(), "handleQueryParams")
	defer span.End()
	perSecondRateParam := queryVals.Get("per_second")
	perMessageSizeParam := queryVals.Get("message_size")
	burstDuration := queryVals.Get("burst_dur")
	s.logger.Info("config params", "perSecondRate", s.config.LogwildPerSecondRate, "outFile", s.config.LogwildOutFile)

	if perSecondRateParam != "" {
		s.logger.Debug("handling perSecondRateParam", "perSecondRateParam", perSecondRateParam)
		perSecondRateInt, err := strconv.ParseInt(perSecondRateParam, 10, 64)
		if err != nil {
			panic(err) // could not convert perSecondRate param to integer :(
		}
		optFuncs = append(optFuncs, logmaker.WithPerSecondRate(perSecondRateInt))
	}
	if perMessageSizeParam != "" {
		s.logger.Debug("handling perMessageSizeParam", "perMessageSizeParam", perMessageSizeParam)
		perMessageSizeInt, err := strconv.ParseInt(perMessageSizeParam, 10, 64)
		if err != nil {
			panic(err) // could not convert perMessageSize param to integer :(
		}
		optFuncs = append(optFuncs, logmaker.WithPerMessageSizeBytes(perMessageSizeInt))
	}
	if burstDuration != "" {
		s.logger.Debug("handling burstDuration", "burstDuration", burstDuration)
		burstDurationInt, err := strconv.ParseInt(burstDuration, 10, 0)
		if err != nil {
			panic(err)
		}
		optFuncs = append(optFuncs, logmaker.WithBurstDuration(time.Second*time.Duration(burstDurationInt)))
	}
	optFuncs = append(optFuncs, logmaker.WithLogger(slog.New(h)))
	s.logger.Info("configured optFuncs", "optFuncs", optFuncs)
	return optFuncs
}

type LogStatsResponse struct {
	logCount int
}
