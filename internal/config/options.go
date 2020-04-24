package config

import (
	"time"

	"github.com/VEuPathDB/lib-go-wdk-api/v0"

	. "github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

type cliOptions struct {
	Auth        string
	Quiet       int
	ReqTimeout  RequestTimeout
	RunSearches bool
	ShowSummary SummaryType
	ThreadNum   uint8
	Verbose     int
	SiteUrl     string
	api         wdk.Api
}

func (c *cliOptions) VerboseLevel() uint8 {
	vb := IntMin(c.Verbose, 2)
	vq := IntMin(c.Quiet, 2)
	return uint8(2 + vb - vq)
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
	return c.SiteUrl
}

func (c *cliOptions) PrintSummary() bool {
	return c.ShowSummary != ""
}

func (c *cliOptions) SummaryType() SummaryType {
	return c.ShowSummary
}

func (c *cliOptions) WdkApi() wdk.Api {
	return c.api
}

func (c *cliOptions) validate() {
	defer x.PanicRecovery()

	c.api = wdk.ForceNew(c.SiteUrl).UseAuthToken(c.Auth)
	c.SiteUrl = c.api.GetUrl().String()
}
