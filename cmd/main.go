package main

import (
	"encoding/json"
	wdk "github.com/VEuPathDB/lib-go-wdk-api/v0"
	"gopkg.in/yaml.v3"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	"github.com/VEuPathDB/script-site-param-cache/internal/config"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"
	"github.com/VEuPathDB/script-site-param-cache/internal/script"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

var version string

func main() {
	defer x.PanicRecovery()
	log.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	start := time.Now()
	opts, validator := config.GetCliOptions(version)

	switch opts.VerboseLevel() {
	case 0: log.SetLevel(log.ErrorLevel)
	case 1: log.SetLevel(log.WarnLevel)
	case 2: log.SetLevel(log.InfoLevel)
	case 3:
		log.SetLevel(log.DebugLevel)
		wdk.Logger().SetLevel(log.DebugLevel)
	case 4:
		log.SetLevel(log.TraceLevel)
		wdk.Logger().SetLevel(log.TraceLevel)
	}
	log.Info("Running param exerciser")

	validator()
	stats := script.NewRunner(opts).Run()

	log.Infof("Completed in %s", time.Now().Sub(start))
	if opts.PrintSummary() {
		printSummary(opts, stats)
	}
}

func printSummary(opts config.CliOptions, stats *out.Summary) {
	stats.Url = opts.BaseUrl()
	switch opts.SummaryType() {
	case "json":
		if err := json.NewEncoder(os.Stdout).Encode(stats); err != nil {
			panic(err)
		}
	case "yaml":
		if err := yaml.NewEncoder(os.Stdout).Encode(stats); err != nil {
			panic(err)
		}
	}
}
