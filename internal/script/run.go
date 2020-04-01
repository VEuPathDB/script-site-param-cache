package script

import (
	"encoding/json"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"os"
	"os/signal"
	"time"

	R "github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

func (r *Runner) Run() {
	recordTypes := make([]string, 0, 25)

	x.FailFast(R.GetRequest(r.url.RecordTypeListUrl()).
		Submit().UnmarshalBody(&recordTypes, R.UnmarshallerFunc(json.Unmarshal)))

	for _, rType := range recordTypes {
		r.processRecordType(rType)
	}

	ch := make(chan os.Signal)
	go func() {
		<- ch
		r.lock.RLock()

		now := time.Now()
		log.InfoFmt("Killing process with %d requests still in progress", len(r.queued))
		for k, v := range r.queued {
			if v.started.Before(v.queued) {
				log.InfoFmt("Queued request to %s. (Queued for %s)", k, now.Sub(v.queued))
			} else {
				log.InfoFmt("Running request to %s. (Running for %s)", k, now.Sub(v.started))
			}
		}
		os.Exit(0)
	}()
	signal.Notify(ch, os.Interrupt)

	for r.wp.WaitingQueueSize() > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	for r.queueLen() > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	r.wp.StopWait()
}

