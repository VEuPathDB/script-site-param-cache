package main

import (
	"time"

	"github.com/VEuPathDB/script-site-param-cache/internal/config"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/script"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

func main() {
	defer x.PanicRecovery()
	start := time.Now()
	opts, validator := config.GetCliOptions()

	log.SetVerbosity(opts.VerboseLevel())
	log.Info("Running script")

	validator()
	script.NewRunner(opts).Run()

	log.InfoFmt("Completed in %s", time.Now().Sub(start))
}
