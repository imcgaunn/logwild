package logmaker

import (
	"fmt"
	"log/slog"
	"time"
)

const (
	microsPerSecond int = 1000 * 1000
	microsPerMilli  int = 1000
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
	tickr := time.NewTicker(tickDuration)
	startTime := time.Now()
	logCount := 0

	lm.Logger.Info("ticker settings", "microsPerEvent", microsPerEvent, "tickDuration", tickDuration, "logsPerTick", logsPerTick, "ticksPerSecond", ticksPerSecond, "logsPerSecond", lm.PerSecondRate)

	// write logsPerTick each tick
	for {
		select {
		case elem := <-tickr.C:
			for i := 0; i < int(logsPerTick); i++ {
				lm.Logger.Debug("processing tick", "elem", elem)
				// actually write the log, and throw up if we can't
				go func() {
					sampleMessage := GetFakeSentence()
					if err := WriteLog(lm, sampleMessage); err != nil {
						panic(err)
					}
					logCount++
				}()
			}

		default:
			// if we have reached our time limit, we can stop. otherwise, just spin on select
			if time.Since(startTime) >= lm.BurstDuration {
				tickr.Stop()
				done <- logCount
				// calculate effective logging rates and return them?
				lm.Logger.Info("completed burst", "logCount", logCount)
				effectiveRateMessages := float64(logCount) / time.Since(startTime).Seconds()
				effectiveRateMbs := (effectiveRateMessages * float64(lm.PerMessageSizeBytes)) / (1024 * 1024)
				lm.Logger.Info(fmt.Sprintf("Effective logging rate: %.2f logs per second\n", effectiveRateMessages))
				lm.Logger.Info(fmt.Sprintf("Effective logging rate (Mb/s): %.2f Mb per second\n", effectiveRateMbs))
				return nil
			}
		}
	}
}

func WriteLog(lm *LogMaker, msg string) error {
	logTime := time.Now().Format(time.RFC3339)
	lm.Logger.Info(msg, "Timestamp", logTime)
	return nil
}
