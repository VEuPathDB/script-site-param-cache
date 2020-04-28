package x

import (
	"context"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

func PanicCatcher(fn func()) func() {
	return func() {
		defer PanicRecovery()
		fn()
	}
}

func PanicRecovery() {
	if rec := recover(); rec != nil {
		if e, ok := rec.(*url.Error); ok {
			if e.Err.Error() == context.DeadlineExceeded.Error() {
				log.Errorf("Timed out while trying to reach %s", e.URL)
				os.Exit(1)
			}
		}

		log.Errorf("%s", rec)
		os.Exit(1)
	}
}
