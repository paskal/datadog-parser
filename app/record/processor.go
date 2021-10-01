package record

import (
	"context"
	"io"
	"time"
)

// Processor goes through records in provided channel and prints alerts and stats on them
type Processor struct {
	LogReader Reader

	records chan []string
}

// Reader is a subset or csv.Reader functions used by the application
type Reader interface {
	Read() ([]string, error)
}

// Start processes new record
// should be called once, is not thread-safe
func (l Processor) Start(ctx context.Context) {
	l.records = make(chan []string)

	// endless log processing loop, wouldn't be terminated using context.
	// leaking goroutine detector (go-leak) would report it.
	go l.readLogRecords()

	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case record := <-l.records:
			l.processRecord(record)
			ticker.Reset(500 * time.Millisecond)
		case <-ticker.C: // ticker ensure we'll recover from alert if no new log entries exist
			l.recalculateAlerts()
		case <-ctx.Done():
			return
		}
	}
}

// processRecord processes new record
func (l Processor) processRecord(record []string) {
}

// recalculateAlerts recalculates alert state
func (l Processor) recalculateAlerts() {
}

// readLogRecords sends new log reader to provided channel
func (l Processor) readLogRecords() {
	for {
		record, err := l.LogReader.Read()
		// ignore all errors, sleep on EOF so that lines could be appended to the log file
		if err == io.EOF {
			// 500ms seems to be a decent compromise between missing not too much data and not burning the CPU away
			time.Sleep(500 * time.Millisecond)
			continue
		}
		l.records <- record
	}
}
