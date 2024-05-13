package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func main() {
	// Configurable parameters
	messageRate := 100 * time.Microsecond // How often to tick (once per ms or 1000 times per second)
	duration := 5 * time.Second           // Total duration to log messages
	ticker := time.NewTicker(messageRate)
	var wg sync.WaitGroup

	// Log file
	file, err := os.Create("log.json")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer file.Close()

	// Start logging in a separate goroutine -- and wait for it to complete
	wg.Add(1)
	go StartLogging(ticker, duration, &wg, file)
	wg.Wait()
}

func StartLogging(tick *time.Ticker, dur time.Duration, wg *sync.WaitGroup, file io.Writer) {
	messageSize := 1024 // Size of each log message in bytes
	sampleMessage := make([]byte, messageSize)
	for i := range sampleMessage {
		sampleMessage[i] = 'A'
	}
	startTime := time.Now()
	logCount := 0
	for elem := range tick.C {
		fmt.Printf("elem %s\n", elem)
		logCount++
		if err := WriteLog(file, sampleMessage); err != nil {
			panic(fmt.Errorf("failed to write log to file"))
		}
		// Check if the logging duration is over
		if time.Since(startTime) >= dur {
			tick.Stop()
			fmt.Printf("Completed. Total logs written: %d\n", logCount)

			// Calculate and print effective logging rate
			effectiveRateMessages := float64(logCount) / time.Since(startTime).Seconds()
			effectiveRateMbs := (effectiveRateMessages * float64(messageSize)) / (1024 * 1024)
			fmt.Printf("Effective logging rate: %.2f logs per second\n", effectiveRateMessages)
			fmt.Printf("Effective logging rate (Mb/s): %.2f Mb per second\n", effectiveRateMbs)
			wg.Done()
		}
	}
}

func WriteLog(w io.Writer, msg []byte) error {
	logEntry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   string(msg),
	}
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}
	w.Write(jsonData)
	w.Write([]byte("\n"))
	return nil
}
