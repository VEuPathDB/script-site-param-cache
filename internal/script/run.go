package script

import (
	"encoding/json"
	"os"
	"os/signal"
	"time"

	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/script-site-param-cache/internal/log"
)

func (r *Runner) Run() {
	recordTypes := make([]string, 0, 25)

	R.GetRequest(r.url.RecordTypeListUrl()).
		SetHttpClient(&r.client).Submit().
		MustUnmarshalBody(&recordTypes, R.UnmarshallerFunc(json.Unmarshal))

	for _, rType := range recordTypes {
		r.processRecordType(rType)
	}

	enableOsSignalHandler(r)

	for r.wp.WaitingQueueSize() > 0 || r.queueLen() > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	r.wp.StopWait()
}

func enableOsSignalHandler(r *Runner) {
	ch := make(chan os.Signal)
	go func() {
		<-ch
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
}
