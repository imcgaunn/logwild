package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func main() {
	// logger
	programLevel := new(slog.LevelVar) // info by default?
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h)) // set up default logger

	// Configurable parameters
	messageRate := 100 * time.Microsecond // How often to tick (10 times per ms or 10000 times per second)
	duration := 5 * time.Second           // Total duration to run program
	ticker := time.NewTicker(messageRate)
	donech := make(chan int) // channel on which signal will be received when logging task is done

	// Start logging in a separate goroutine -- and wait for it to complete
	slog.Info("about to start logging with the following settings", "messageRate", messageRate, "duration", duration)
	go StartLogging(ticker, duration, donech)
	totalLogs := <-donech
	slog.Info("program done", "totalLogs", totalLogs)
}

func StartLogging(tick *time.Ticker, dur time.Duration, donech chan int) {
	messageSize := 1024 // Size of each log message in bytes
	sampleMessage := make([]byte, messageSize)
	for i := range sampleMessage {
		sampleMessage[i] = 'A'
	}
	startTime := time.Now()
	logCount := 0
	for {
		select {
		case elem := <-tick.C:
			slog.Info("processing tick", "elem", elem)
			logCount++
			if err := WriteLog(sampleMessage); err != nil {
				panic(err)
			}
		default:
			// Check if the logging duration is over
			if time.Since(startTime) >= dur {
				tick.Stop()
				donech <- logCount
				slog.Info("Completed", "logCount", logCount)

				// Calculate and print effective logging rate
				effectiveRateMessages := float64(logCount) / time.Since(startTime).Seconds()
				effectiveRateMbs := (effectiveRateMessages * float64(messageSize)) / (1024 * 1024)
				fmt.Printf("Effective logging rate: %.2f logs per second\n", effectiveRateMessages)
				fmt.Printf("Effective logging rate (Mb/s): %.2f Mb per second\n", effectiveRateMbs)
				return
			}
		}
	}
}

func WriteLog(msg []byte) error {
	logEntry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   string(msg),
	}
	slog.Info(logEntry.Message, "Timestamp", logEntry.Timestamp)
	return nil
}
