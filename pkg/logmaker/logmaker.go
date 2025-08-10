package logmaker

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	microsPerSecond int = 1000 * 1000
	microsPerMilli  int = 1000
)

type OptFunc func(*Opts)

type Opts struct {
	PerSecondRate  int64
	PerMessageSize int64
	BurstDuration  time.Duration
	Logger         *slog.Logger
}

type LogMaker struct {
	Opts
}

func defaultOpts() Opts {
	return Opts{
		PerSecondRate:  1000,
		PerMessageSize: 48,
		BurstDuration:  5 * time.Second,
		Logger:         slog.Default(),
	}
}

func WithPerSecondRate(psr int64) OptFunc {
	return func(opts *Opts) {
		opts.PerSecondRate = psr
	}
}

func WithPerMessageSize(b int64) OptFunc {
	return func(opts *Opts) {
		opts.PerMessageSize = b
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
	microsPerEvent := float64(microsPerSecond) / float64(lm.PerSecondRate)
	lm.Logger.Debug("about to start logMaker", "microsPerEvent", microsPerEvent, "perSecondRate", lm.PerSecondRate)
	logsPerTick := float64(1) // how many logs need to be emitted per tick
	var ticksPerSecond int64
	var tickDuration time.Duration
	if microsPerEvent < 5*float64(microsPerMilli) {
		// if each event is < 5ms we have to do some tricks
		// maximum resolution of go ticker is about 1ms on unix
		lm.Logger.Debug("microsPerEvent less than microsPerMilli", "microsPerMilli", microsPerMilli)
		tickDuration = time.Duration(5 * time.Millisecond)
	} else {
		// if time between events is more than 1ms, don't worry about, it just
		// do one log per tick.
		tickDuration = time.Duration(microsPerEvent) * time.Microsecond
	}
	ticksPerSecond = int64(time.Second / tickDuration)
	logsPerTick = float64(lm.PerSecondRate) / float64(ticksPerSecond)

	// --- Worker Pool Setup ---
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, lm.PerSecondRate)
	var wg sync.WaitGroup
	var logCount int64

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range jobs {
				if err := WriteLog(lm, msg); err != nil {
					lm.Logger.Error("failed to write log", "error", err)
				} else {
					atomic.AddInt64(&logCount, 1)
				}
			}
		}()
	}
	// --- End Worker Pool Setup ---

	tickr := time.NewTicker(tickDuration)
	defer tickr.Stop() // Ensure ticker is stopped

	burstDone := time.After(lm.BurstDuration)
	startTime := time.Now()

	lm.Logger.Info("ticker settings", "microsPerEvent", microsPerEvent, "tickDuration", tickDuration, "logsPerTick", logsPerTick, "ticksPerSecond", ticksPerSecond, "logsPerSecond", lm.PerSecondRate, "numWorkers", numWorkers)

	// Main loop
	for {
		select {
		case <-burstDone:
			close(jobs) // Signal workers to stop
			wg.Wait()   // Wait for workers to finish
			finalLogCount := atomic.LoadInt64(&logCount)
			done <- int(finalLogCount)

			// calculate effective logging rates and return them?
			timeSpentSeconds := time.Since(startTime).Seconds()
			lm.Logger.Info("completed burst",
				"timeSpentSeconds", fmt.Sprintf("%.2f", timeSpentSeconds),
				"logCount", finalLogCount)
			effectiveRateMessages := float64(finalLogCount) / time.Since(startTime).Seconds()
			effectiveRateMbs := (effectiveRateMessages * float64(lm.PerMessageSize) * 9) / (1024 * 1024)
			lm.Logger.Info(fmt.Sprintf("Effective logging rate: %.2f logs per second", effectiveRateMessages))
			lm.Logger.Info(fmt.Sprintf("Effective logging rate (Mb/s): %.2f Mb per second", effectiveRateMbs))
			return nil
		case <-tickr.C:
			for i := 0; i < int(logsPerTick); i++ {
				sampleMessage := GetFakeSentence(int(lm.PerMessageSize))
				// Use a non-blocking send in case the workers are backed up
				select {
				case jobs <- sampleMessage:
				default:
					lm.Logger.Warn("log message dropped, channel full")
				}
			}
		}
	}
}

func WriteLog(lm *LogMaker, msg string) error {
	logTime := time.Now().Format(time.RFC3339)
	lm.Logger.Info(msg, "Timestamp", logTime)
	return nil
}
