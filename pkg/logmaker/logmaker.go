package logmaker

import (
	"fmt"
	"log/slog"
	"time"
)

type OptFunc func(*Opts)

type Opts struct {
	PerSecondRate       int64
	PerMessageSizeBytes int64
	BurstDuration       time.Duration
	Logger              *slog.Logger
}

type LogMaker struct {
	Opts
}

func defaultOpts() Opts {
	return Opts{
		PerSecondRate:       1000,
		PerMessageSizeBytes: 1024 * 2,
		BurstDuration:       5 * time.Second,
		Logger:              slog.Default(),
	}
}

func WithPerSecondRate(psr int64) OptFunc {
	return func(opts *Opts) {
		opts.PerSecondRate = psr
	}
}

func WithPerMessageSizeBytes(b int64) OptFunc {
	return func(opts *Opts) {
		opts.PerMessageSizeBytes = b
	}
}

func WithBurstDuration(d time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.BurstDuration = d
	}
}

func WithLogger(l *slog.Logger) OptFunc {
	return func(opts *Opts) {
		opts.Logger = l
	}
}

func NewLogMaker(opts ...OptFunc) *LogMaker {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}
	return &LogMaker{o}
}

func (lm *LogMaker) StartWriting(done chan int) error {
	// calculate duration based on PerSecondRate in cfg
	// just always use microsecond precision
	// microseconds between ticks
	microsPerSecond := time.Second / time.Microsecond
	microsPerEvent := float64(microsPerSecond) / float64(lm.PerSecondRate)
	tickDuration := time.Duration(microsPerEvent) * time.Microsecond
	tickr := time.NewTicker(tickDuration)
	startTime := time.Now()
	logCount := 0
	sampleMessage := make([]byte, lm.PerMessageSizeBytes)
	for i := range sampleMessage {
		sampleMessage[i] = 'A'
	}
	for {
		select {
		case elem := <-tickr.C:
			lm.Logger.Debug("processing tick", "elem", elem)
			// actually write the log, and throw up if we can't
			if err := WriteLog(lm, sampleMessage); err != nil {
				panic(err)
			}
			logCount++
		default:
			// if we have reached our time limit, we can stop. otherwise, just spin on select
			if time.Since(startTime) >= lm.BurstDuration {
				tickr.Stop()
				done <- logCount
				lm.Logger.Debug("completed burst", "logCount", logCount)
				// calculate effective logging rates and return them?
				effectiveRateMessages := float64(logCount) / time.Since(startTime).Seconds()
				effectiveRateMbs := (effectiveRateMessages * float64(lm.PerMessageSizeBytes)) / (1024 * 1024)
				fmt.Printf("Effective logging rate: %.2f logs per second\n", effectiveRateMessages)
				fmt.Printf("Effective logging rate (Mb/s): %.2f Mb per second\n", effectiveRateMbs)
				return nil
			}
		}
	}
}

func WriteLog(lm *LogMaker, msg []byte) error {
	logTime := time.Now().Format(time.RFC3339)
	lm.Logger.Info(string(msg), "Timestamp", logTime)
	return nil
}
