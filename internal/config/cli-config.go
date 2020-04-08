package config

import (
	"net/http"
	"strings"
	"time"

	. "github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

var staticOptions *cliOptions

func init() {
	staticOptions = parseCliOptions()
}

type CliOptions interface {
	VerboseLevel() uint8
	Threads() uint8
	AuthToken() string
	SearchEnabled() bool
	RequestTimeout() time.Duration
	BaseUrl() string
}

type ValidatorFunc func()

func GetCliOptions() (CliOptions, ValidatorFunc) {
	return staticOptions, staticOptions.validate
}

type cliOptions struct {
	Verbose []bool `short:"v" long:"verbose" description:"Enable verbose log output. Can be specified a second time for more verbose logging"`

	ThreadNum uint8 `short:"p" long:"parallel" default:"16" description:"Number of threads to run on"`

	Auth string `long:"auth" value-name:"auth_tkt" description:"QA Auth Token.\nThis can be retrieved by logging in to a QA site and pulling the value from either the \"auth_tkt\" query parameter or the cookie with same name."`

	RunSearches bool `short:"r" long:"run-searches" description:"Set to attempt to run all the searches"`

	ReqTimeout RequestTimeout `short:"t" long:"timeout" default:"10m" description:"Max duration cap on individual requests.\nFormatted as <num><unit>[<num><unit>...] for example \"5m\" for five minutes or \"2m30s\" for two minutes and thirty seconds.\n\nValid units are:\n  ms = milliseconds\n  s  = seconds\n  m  = minutes\n  h  = hours\n"`

	Positional struct {
		Url string `positional-arg-name:"URL" description:"Site URL\nExample: https://plasmodb.org"`
	} `positional-args:"yes" required:"1"`
}

func (c *cliOptions) VerboseLevel() uint8 {
	return uint8(len(c.Verbose))
}

func (c *cliOptions) Threads() uint8 {
	return c.ThreadNum
}

func (c *cliOptions) AuthToken() string {
	return c.Auth
}

func (c *cliOptions) SearchEnabled() bool {
	return c.RunSearches
}

func (c *cliOptions) RequestTimeout() time.Duration {
	return c.ReqTimeout.ToDuration()
}

func (c *cliOptions) BaseUrl() string {
	return c.Positional.Url
}

func (c *cliOptions) validate() {
	defer x.PanicRecovery()

	if c.ThreadNum < 1 || c.ThreadNum > 24 {
		panic("Invalid number of threads: '%d'.  Must be in the range [1..24].")
	}

	res := GetRequest(c.Positional.Url, &http.Client{
		Timeout: c.ReqTimeout.ToDuration(),
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}})

	if res.MustGetResponseCode() == http.StatusFound {
		c.Positional.Url = res.MustGetHeader("Location")
	}

	if !strings.HasSuffix(c.Positional.Url, "/") {
		c.Positional.Url = c.Positional.Url + "/"
	}
}

func parseCliOptions() (opts *cliOptions) {
	opts = new(cliOptions)
	_ = x.ParseFlags(opts)
	return
}
