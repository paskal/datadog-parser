package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
)

type opts struct {
	FilePath                string        `long:"filepath" env:"FILEPATH" default:"" description:"csv file path, stdin is used if not specified"`
	AlertWindow             time.Duration `long:"alert_window" env:"ALERT_WINDOW" default:"2m" description:"alert windows"`
	AlertThresholdPerSecond int           `long:"alert_threshold_per_sec" env:"ALERT_THRESHOLD_PER_SEC" default:"10" description:"threshold for alert, requests per second"`
}

// recordReader is a subset or csv.Reader functions used by the application
type recordReader interface {
	Read() ([]string, error)
}

func main() {
	var opts opts
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("Unable to parse the args: %v", err)
		os.Exit(2)
	}

	var logReader recordReader
	var err error

	// retrieve the recordReader either from file or from stdin
	if opts.FilePath != "" {
		var file io.Closer
		logReader, file, err = getFileReader(opts.FilePath)
		if err != nil {
			log.Printf("Error opening csv file: %v", err)
			os.Exit(3)
		}
		defer file.Close()
	} else {
		logReader = getStdinReader()
	}

	// catch TERM signal and invoke graceful termination
	// otherwise it's impossible to test main()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	recordsChan := make(chan []string)
	// endless log processing loop, wouldn't be terminated using context
	go readLogRecords(recordsChan, logReader)

	ticker := time.NewTicker(500 * time.Millisecond)
	// main logic loop
	for {
		select {
		case record := <-recordsChan:
			processRecord(record)
			ticker.Reset(500 * time.Millisecond)
		case <-ticker.C: // ticker ensure we'll recover from alert if no new log entries exist
			recalculateAlerts()
		case <-ctx.Done():
			return
		}
	}
}

// processRecord processes new record
func processRecord(record []string) {
}

// recalculateAlerts recalculates alert state
func recalculateAlerts() {
}

// getFileReader returns recordReader from given file
func getFileReader(filePath string) (recordReader, io.Closer, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening %v: %v", filePath, err)
	}
	return csv.NewReader(f), f, nil
}

// getStdinReader returns recordReader from stdin
func getStdinReader() recordReader {
	return csv.NewReader(os.Stdin)
}

// readLogRecords sends new log reader to provided channel
func readLogRecords(recordsChan chan<- []string, logReader recordReader) {
	for {
		record, err := logReader.Read()
		// ignore all errors, sleep on EOF so that lines could be appended to the log file
		if err == io.EOF {
			// 500ms seems to be a decent compromise between missing not too much data and not burning the CPU away
			time.Sleep(500 * time.Millisecond)
			continue
		}
		recordsChan <- record
	}
}
