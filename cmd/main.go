package main

import (
	"encoding/json"
	wdk "github.com/VEuPathDB/lib-go-wdk-api/v0"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"

	"github.com/VEuPathDB/script-site-param-cache/internal/config"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"
	"github.com/VEuPathDB/script-site-param-cache/internal/script"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

var version string

func main() {
	defer x.PanicRecovery()
	start := time.Now()
	wdk.Logger().SetLevel(logrus.TraceLevel)
	opts, validator := config.GetCliOptions(version)

	log.SetVerbosity(opts.VerboseLevel())
	log.Info("Running param exerciser")

	validator()
	stats := script.NewRunner(opts).Run()

	log.InfoFmt("Completed in %s", time.Now().Sub(start))
	if opts.PrintSummary() {
		printSummary(opts, stats)
	}
}

func printSummary(opts config.CliOptions, stats out.Summary) {
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
