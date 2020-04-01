package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Foxcapades/Go-ChainRequest"
	req "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

const (
	testTimeout time.Duration = 3 * time.Second
)

type CliOptions struct {
	Verbose bool `short:"v" long:"verbose" description:"Enable verbose log output"`

	Threads uint8 `short:"t" long:"threads" default:"16" description:"Number of threads to run on"`

	Auth string `long:"auth" value-name:"auth_tkt" description:"QA Auth Token"`

	Positional struct {
		Url  string `positional-arg-name:"URL" description:"Site URL\nExample: https://plasmodb.org"`
	} `positional-args:"yes" required:"1"`
}

func (c *CliOptions) Validate() {
	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(*url.Error); ok {
				if e.Err.Error() == context.DeadlineExceeded.Error() {
					panic(
						fmt.Sprintf("Could not connect to site %s within the timeout of %s",
							c.Positional.Url, testTimeout))
				}
			}
			panic(rec)
		}
	}()

	if c.Threads < 1 || c.Threads > 16 {
		panic("Invalid number of threads: '%d'.  Must be in the range [1..16].")
	}

	res := req.GetRequest(c.Positional.Url).
		DisableRedirects().
		SetRequestBuilder(creq.RequestBuilderFunc(func(r creq.Request) (*http.Request, error) {
			request, err := http.NewRequest(string(r.GetMethod()), r.GetUrl(), nil)
			x.FailFast(err)

			ctx, _ := context.WithTimeout(context.Background(), testTimeout)

			return request.WithContext(ctx), nil
		})).
		Submit()

	x.FailFast(res.GetError())

	if res.MustGetResponseCode() == http.StatusFound {
		c.Positional.Url = res.MustGetHeader("Location")
	}

	if !strings.HasSuffix(c.Positional.Url, "/") {
		c.Positional.Url = c.Positional.Url + "/"
	}

}

func ParseCliOptions() (opts *CliOptions) {
	opts = new(CliOptions)
	_ = x.ParseFlags(opts)
	return
}
