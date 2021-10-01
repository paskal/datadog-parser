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

	"github.com/paskal/datadog-parser/app/record"
)

type opts struct {
	FilePath                string        `long:"filepath" env:"FILEPATH" default:"" description:"csv file path, stdin is used if not specified"`
	AlertWindow             time.Duration `long:"alert_window" env:"ALERT_WINDOW" default:"2m" description:"alert windows"`
	AlertThresholdPerSecond int           `long:"alert_threshold_per_sec" env:"ALERT_THRESHOLD_PER_SEC" default:"10" description:"threshold for alert, requests per second"`
}

func main() {
	var opts opts
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("Unable to parse the args: %v", err)
		os.Exit(2)
	}

	var logReader record.Reader
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

	logProcessor := record.Processor{LogReader: logReader}
	logProcessor.Start(ctx)
}

// getFileReader returns recordReader from given file
func getFileReader(filePath string) (record.Reader, io.Closer, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening %v: %v", filePath, err)
	}
	return csv.NewReader(f), f, nil
}

// getStdinReader returns recordReader from stdin
func getStdinReader() record.Reader {
	return csv.NewReader(os.Stdin)
}
