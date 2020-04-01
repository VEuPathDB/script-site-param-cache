package main

import (
	"os"
	"time"

	"github.com/VEuPathDB/script-site-param-cache/internal/config"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/script"
)

func main() {
	defer recov()
	start := time.Now()
	log.Info("Parsing & Validating Configuration")

	opts := config.ParseCliOptions()
	opts.Validate()
	log.ConfigureLogger(opts)

	log.Info("Running script")
	script.NewRunner(opts).Run()
	log.InfoFmt("Completed in %s", time.Now().Sub(start))
}

func recov() {
	if rec := recover(); rec != nil {
		log.ErrorFmt("%s", rec)
		os.Exit(1)
	}
}
