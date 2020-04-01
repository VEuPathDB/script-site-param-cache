package x

import (
	"github.com/jessevdk/go-flags"
	"os"
)

func ParseFlags(val interface{}) []string {
	out, err := flags.Parse(val)

	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		os.Exit(0)
	}

	FailFast(err)

	return out
}


func FailFast(err error) {
	if err != nil {
		panic(err)
	}
}