package main

import (
	"context"
	"encoding/csv"
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

	if opts.AlertWindow == 0 {
		log.Print("Alert window must be non-zero")
		os.Exit(2)
	}

	if opts.AlertThresholdPerSecond == 0 {
		log.Print("Alert threshold must be non-zero")
		os.Exit(2)
	}

	var logReader *csv.Reader

	// retrieve the recordReader either from file or from stdin
	if opts.FilePath != "" {
		f, err := os.Open(opts.FilePath)
		if err != nil {
			log.Printf("Error opening csv file: %v", err)
			os.Exit(3)
		}
		defer f.Close()
		logReader = csv.NewReader(f)
	} else {
		logReader = csv.NewReader(os.Stdin)
	}
	logReader.FieldsPerRecord = 7

	// catch TERM signal and invoke graceful termination
	// otherwise it's impossible to test main()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	logProcessor := record.Processor{
		LogReader:               logReader,
		AlertWindow:             opts.AlertWindow,
		AlertThresholdPerSecond: opts.AlertThresholdPerSecond,
	}
	logProcessor.Start(ctx)
}
