package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Foxcapades/Argonaut/v0"
	"github.com/Foxcapades/Argonaut/v0/pkg/argo"
)

type CliOptions interface {
	VerboseLevel() uint8
	Threads() uint8
	AuthToken() string
	SearchEnabled() bool
	RequestTimeout() time.Duration
	BaseUrl() string
	PrintSummary() bool
	SummaryType() SummaryType
}

type ValidatorFunc func()

const (
	fHelpVerb = "Enable verbose log output.  Can be specified a second time for" +
		" more verbose logging"
	fHelpQuiet = "Print less output.  Specified once will disabled info level " +
		"logging, specified twice will disable all logging."
	fHelpThreads = "Number of threads on which to run concurrent requests"
	fHelpAuth = "QA Auth Token.\nThis can be retrieved by logging in to a QA " +
		"site and pulling the value from either the \"auth_tkt\" query parameter " +
		"or the cookie with same name."
	fHelpRun  = "Set to attempt to run all the searches"
	fHelpTime = "Max duration cap on individual requests.\nFormatted as " +
		"<num><unit>[<num><unit>...] for example \"5m\" for five minutes or " +
		"\"2m30s\" for two minutes and thirty seconds.\n\nValid units are:\n  " +
		"ms = milliseconds\n  s  = seconds\n  m  = minutes\n  h  = hours\n"
	fHelpSummary = "Print a result summary of requests that passed vs failed." +
		"\nThis output is not affected by the --quiet flag.\n\nValid Types:\n  " +
		"json\n  yaml"
	fHelpVersion = "Prints the script version"
)

func GetCliOptions(version string) (CliOptions, ValidatorFunc) {
	out := cliOptions{}

	cli.NewCommand().
		Flag(cli.LFlag("auth", fHelpAuth).
			Arg(cli.NewArg().Name("auth_tkt").Bind(&out.Auth).Require())).
		Flag(cli.SlFlag('p', "parallel", fHelpThreads).
			Arg(cli.NewArg().Default(uint8(16)).Bind(&out.ThreadNum).Require())).
		Flag(cli.SlFlag('q', "quiet", fHelpQuiet).BindUseCount(&out.Quiet)).
		Flag(cli.LFlag("run-searches", fHelpRun).
			Bind(&out.RunSearches, false)).
		Flag(cli.SlFlag('t', "timeout", fHelpTime).
			Arg(cli.NewArg().
				Name("timeout").
				Default("10m").
				Bind(&out.ReqTimeout).
				Require())).
		Flag(cli.SlFlag('v', "verbose", fHelpVerb).BindUseCount(&out.Verbose)).
		Flag(cli.LFlag("summary", fHelpSummary).
			Arg(cli.NewArg().Name("json|yaml").Bind(&out.ShowSummary).Require())).
		Flag(cli.SlFlag('V', "version", fHelpVersion).
			OnHit(func(argo.Flag) {
				fmt.Println(version)
				os.Exit(0)
			})).
		Arg(cli.NewArg().Name("site-url").Bind(&out.SiteUrl).Require()).
		MustParse()

	return &out, out.validate
}
