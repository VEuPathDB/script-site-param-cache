package x

import (
	"context"
	"net/url"
	"os"

	"github.com/VEuPathDB/script-site-param-cache/internal/log"
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
				log.ErrorFmt("Timed out while trying to reach %s", e.URL)
				os.Exit(1)
			}
		}

		log.ErrorFmt("%s", rec)
		os.Exit(1)
	}
}
