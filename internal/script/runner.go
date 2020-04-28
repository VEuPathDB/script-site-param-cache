package script

import (
	"net/http"
	"sync"
	"time"

	"github.com/VEuPathDB/lib-go-wdk-api/v0"
	"github.com/gammazero/workerpool"
	log "github.com/sirupsen/logrus"

	"github.com/VEuPathDB/lib-go-rest-types/veupath"
	"github.com/VEuPathDB/script-site-param-cache/internal/config"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"
)

type status struct {
	queued  time.Time
	started time.Time
}

type Runner struct {
	stats  out.Summary
	queued map[string]*status
	lock   sync.RWMutex
	wp     *workerpool.WorkerPool
	url    veupath.ApiUrlBuilder
	opts   config.CliOptions
	client http.Client
	api    wdk.Api
}

func NewRunner(opt config.CliOptions) (runner *Runner) {
	runner = &Runner{
		queued: make(map[string]*status),
		wp:     workerpool.New(int(opt.Threads())),
		url:    veupath.NewApiUrlBuilder(opt.BaseUrl()),
		opts:   opt,
		client: http.Client{Timeout: opt.RequestTimeout()},
		api:    opt.WdkApi(),
	}
	runner.url.SetAuthTkt(opt.AuthToken())
	printSetup(runner)
	return
}

func (r *Runner) push(url string) {
	r.lock.Lock()
	if _, ok := r.queued[url]; ok {
		panic("More than one simultaneous request to the same url: " + url)
	}
	r.queued[url] = &status{queued: time.Now()}
	r.lock.Unlock()
}

func (r *Runner) start(url string) {
	r.lock.RLock()
	if _, ok := r.queued[url]; !ok {
		panic("Attempted to start a non-queued request to: " + url)
	}
	r.queued[url].started = time.Now()
	r.lock.RUnlock()
}

func (r *Runner) pop(url string) {
	r.lock.Lock()
	delete(r.queued, url)
	r.lock.Unlock()
}

func (r *Runner) queueLen() (out int) {
	r.lock.RLock()
	out = len(r.queued)
	r.lock.RUnlock()
	return
}

func printSetup(r *Runner) {
	log.Debug("Setting up runner")
	log.Tracef(
		`Base URL:     %s
     Timeout:      %s
     Threads:      %d
     Run Searches: %t`,
		r.opts.BaseUrl(),
		r.opts.RequestTimeout(),
		r.opts.Threads(),
		r.opts.SearchEnabled(),
	)
}
