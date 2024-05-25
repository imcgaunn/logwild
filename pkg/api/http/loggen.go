package http

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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
	s.logger.Info("lm config", "perSecondRate", lm.PerSecondRate)
	donech := make(chan int)
	defer close(donech)
	go func() {
		err := lm.StartWriting(donech)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			s.logger.Error("encountered problem trying to start logmaker", "err", err)
			// should return bad status to http as well
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
	_, span := s.tracer.Start(r.Context(), "handleQueryParams")
	defer span.End()

	s.logger.Info("config params",
		"perSecondRate", s.config.LogwildPerSecondRate,
		"outFile", s.config.LogwildOutFile)

	perSecondRateInt, err := s.tryParseAndLogIntParam(r, "per_second")
	if err == nil {
		optFuncs = append(optFuncs, logmaker.WithPerSecondRate(perSecondRateInt))
	}
	perMessageSizeInt, err := s.tryParseAndLogIntParam(r, "message_size")
	if err == nil {
		optFuncs = append(optFuncs, logmaker.WithPerMessageSize(perMessageSizeInt))
	}
	burstDurationInt, err := s.tryParseAndLogIntParam(r, "burst_dur")
	if err == nil {
		optFuncs = append(optFuncs, logmaker.WithBurstDuration(time.Duration(burstDurationInt)*time.Second))
	}
	optFuncs = append(optFuncs, logmaker.WithLogger(slog.New(h)))
	s.logger.Info("configured optFuncs", "optFuncs", optFuncs)
	return optFuncs
}

func (s *Server) tryParseAndLogIntParam(r *http.Request, paramName string) (int64, error) {
	queryVals := r.URL.Query()
	paramVal := queryVals.Get(paramName)
	s.logger.Debug("handling parameter", "paramName", paramName, "paramVal", paramVal)
	if paramVal == "" {
		return 0, errors.New("requested parameter not present in request")
	}
	intValue, err := strconv.ParseInt(paramVal, 10, 64)
	if err != nil {
		s.logger.Error("could not parse param as integer", "paramName", paramName, "paramVal", paramVal, "err", err)
		return 0, errors.New("could not parse param as integer")
	}
	return intValue, nil
}

type LogStatsResponse struct {
	logCount int
}
