package script

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"

	"github.com/VEuPathDB/script-site-param-cache/internal/out"
)

// Run executes the param cache script logic.
//
// It first retrieves the list of record types, then for
// each record type name, it looks up the record type to get
// that record type's searches.  Then for each search on
// each record type retrieves the full search json for the
// search name
//
//     record_names = get_record_names()
//     foreach(record_name in record_names)
//       record_details = get_record_details(record_name)
//       foreach(search_summary in record_details.searches)
//         search_details = get_search_details(search_summary)
func (r *Runner) Run() *out.Summary {
	start := time.Now()
	recordTypes := r.api.MustGetExpandedRecordTypes()
	log.WithField("time", time.Since(start)).Debug("Successful record-types GET")

	for i := range recordTypes {
		rt := &recordTypes[i]
		for j := range rt.Searches {
			r.processShortSearch(rt, &rt.Searches[j])
		}
	}

	enableOsSignalHandler(r)

	for r.wp.WaitingQueueSize() > 0 || r.queueLen() > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	r.wp.StopWait()
	return &r.stats
}

func enableOsSignalHandler(r *Runner) {
	ch := make(chan os.Signal)
	go func() {
		<-ch
		r.lock.RLock()

		now := time.Now()
		log.Infof("Killing process with %d requests still in progress", len(r.queued))
		for k, v := range r.queued {
			if v.started.Before(v.queued) {
				log.Infof("Queued request to %s. (Queued for %s)", k, now.Sub(v.queued))
			} else {
				log.Infof("Running request to %s. (Running for %s)", k, now.Sub(v.started))
			}
		}
		os.Exit(0)
	}()
	signal.Notify(ch, os.Interrupt)
}
