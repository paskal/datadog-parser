package record

import (
	"context"
	"fmt"
	"io"
	"sort"
	"time"
)

const reportInterval = time.Second * 10

var printFunction = fmt.Printf // overwritten in tests

// Reader is a subset or csv.Reader functions used by the application
type Reader interface {
	Read() ([]string, error)
}

// Processor goes through records in provided channel and prints alerts and stats on them
type Processor struct {
	LogReader               Reader
	AlertWindow             time.Duration
	AlertThresholdPerSecond int

	alertState     bool
	lastReport     time.Time
	records        chan []string
	history        map[int64]historyRecord
	historicalHits int
}

// Start processes new record from provided Reader
// should be called once, is not thread-safe
func (l *Processor) Start(ctx context.Context) {
	l.records = make(chan []string)
	l.history = make(map[int64]historyRecord)

	go l.readLogRecords(ctx)

	for {
		select {
		case record := <-l.records:
			l.processRecord(record)
		case <-ctx.Done():
			return
		}
	}
}

// processRecord processes new record
func (l *Processor) processRecord(rawRecord []string) {
	r := parseRecord(rawRecord)
	if r == nil {
		return
	}

	l.historicalHits++
	ts := r.date.Unix()
	history, ok := l.history[ts]
	if !ok {
		history = newHistoryRecord()
	}
	history.add(r)
	l.history[ts] = history

	if l.lastReport.Equal(time.Time{}) {
		l.lastReport = r.date
	}

	if r.date.Sub(l.lastReport) >= reportInterval {
		// we need to print the report on entries before the last one, as new log entry might be hours away from previous one
		l.printReport(l.findPreLastReport(r.date))
		l.lastReport = r.date
	}
	l.cleanHistory(r.date)
	l.recalculateAlerts(r.date)
}

// findPreLastReport finds the date of the last report before the provided time
func (l *Processor) findPreLastReport(lastEntry time.Time) time.Time {
	preLastEntryTime := time.Time{}
	for k := range l.history {
		t := time.Unix(k, 0)
		if t.Before(lastEntry) && !t.Equal(lastEntry) && preLastEntryTime.Before(t) {
			preLastEntryTime = t
		}
	}
	return preLastEntryTime
}

// printReport for the reportInterval
func (l *Processor) printReport(lastEntry time.Time) {
	stats := newHistoryRecord()
	for k, v := range l.history {
		// we check if the entry is before the last one to prevent the very last log entry from being counted
		if !time.Unix(k, 0).After(lastEntry) && lastEntry.Sub(time.Unix(k, 0)) <= reportInterval {
			stats.append(v)
		}
	}
	var topSection string
	var topHits int

	// sort sections so that they appear in the output in the same order reliably
	sections := []string{}
	for k := range stats.sections {
		sections = append(sections, k)
	}
	sort.Strings(sections)

	for _, k := range sections {
		if topHits == stats.sections[k] {
			topSection += " and " + k
		}
		if topHits < stats.sections[k] {
			topHits = stats.sections[k]
			topSection = k
		}
	}

	printFunction("%s: %d hits from %d users with %d bytes transferred, top %s with %d hits\n", //nolint:errcheck
		lastEntry.In(time.UTC),
		stats.hits,
		len(stats.uniqueUsers),
		stats.bytesTransferred,
		topSection,
		topHits,
	)
}

// cleanHistory drops history older than AlertWindow from specified date
func (l *Processor) cleanHistory(d time.Time) {
	for k := range l.history {
		if d.Sub(time.Unix(k, 0)) > l.AlertWindow {
			l.historicalHits -= l.history[k].hits
			delete(l.history, k)
		}
	}
}

// recalculateAlerts recalculates alert state
func (l *Processor) recalculateAlerts(currentTime time.Time) {
	hitsPerSecond := float64(l.historicalHits) / l.AlertWindow.Seconds()

	if hitsPerSecond > float64(l.AlertThresholdPerSecond) {
		if !l.alertState {
			printFunction("%s: Alert RED, ~%.2f hits per second which is higher than %d (%d total) in the last %s\n", //nolint:errcheck
				currentTime.In(time.UTC),
				hitsPerSecond,
				l.AlertThresholdPerSecond,
				l.historicalHits,
				l.AlertWindow,
			)
			l.alertState = true
		}
		return
	}

	if l.alertState {
		printFunction("%s: Alert GREEN, ~%.2f hits per second which is lower than %d (%d total) in the last %s\n", //nolint:errcheck
			currentTime.In(time.UTC),
			hitsPerSecond,
			l.AlertThresholdPerSecond,
			l.historicalHits,
			l.AlertWindow,
		)
		l.alertState = false
	}
}

// readLogRecords sends new log reader to provided channel,
// wouldn't be terminated using context unless there is new log entry,
// but would reliably terminate in tests with properly constructed Reader
func (l *Processor) readLogRecords(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
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

type historyRecord struct { // key is unix timestamp
	bytesTransferred int                 // used for stats
	hits             int                 // used for alerting
	sections         map[string]int      // hit stats per section
	uniqueUsers      map[string]struct{} // unique user counter
}

func newHistoryRecord() historyRecord {
	return historyRecord{
		sections:    make(map[string]int),
		uniqueUsers: make(map[string]struct{}),
	}
}

func (h *historyRecord) add(r *record) {
	h.bytesTransferred += r.bytes
	h.sections[r.section]++
	h.uniqueUsers[r.remotehost] = struct{}{}
	h.hits++
}

func (h *historyRecord) append(new historyRecord) {
	h.bytesTransferred += new.bytesTransferred
	h.hits += new.hits
	for s := range new.sections {
		h.sections[s]++
	}
	for u := range new.uniqueUsers {
		h.uniqueUsers[u] = struct{}{}
	}
}
