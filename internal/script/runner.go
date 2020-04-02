package script

import (
	"github.com/gammazero/workerpool"
	"net/http"
	"sync"
	"time"

	"github.com/VEuPathDB/lib-go-rest-types/veupath"
	"github.com/VEuPathDB/script-site-param-cache/internal/config"
)

type status struct {
	queued  time.Time
	started time.Time
}

type Runner struct {
	queued map[string]*status
	lock   sync.RWMutex
	wp     *workerpool.WorkerPool
	url    veupath.ApiUrlBuilder
	opts   *config.CliOptions
	client http.Client
}

func NewRunner(opt *config.CliOptions) (runner *Runner) {
	runner = &Runner{
		queued: make(map[string]*status),
		wp:     workerpool.New(int(opt.Threads)),
		url:    veupath.NewApiUrlBuilder(opt.Positional.Url),
		opts:   opt,
		client: http.Client{Timeout: time.Duration(opt.RequestTimeout)},
	}
	runner.url.SetAuthTkt(opt.Auth)
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