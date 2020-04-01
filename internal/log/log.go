package log

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/VEuPathDB/script-site-param-cache/internal/config"
)

var verbose bool

const (
	prefixError = "\033[91mERROR\033[0m"
	prefixWarn  = "\033[33mWARN \033[0m"
	prefixInfo  = "\033[32mINFO \033[0m"
	prefixDebug = "\033[36mDEBUG\033[0m"
)

const (
	timeStampFmt = "[2006-01-02T15:04:05.000000000Z07:00]"
)

var nlPadding string
var replace = regexp.MustCompile("\n[ \t]*")

func init() {
	buf := strings.Builder{}
	ln  := len(timeStampFmt) + 8

	buf.Grow(ln)
	buf.Reset()

	buf.WriteByte('\n')
	for i := 1; i < ln; i++ {
		buf.WriteByte(' ')
	}

	nlPadding = buf.String()
}

func nowStamp() string {
	return time.Now().Format(timeStampFmt)
}

func nlPad(val string) string {
	return replace.ReplaceAllString(val, nlPadding)
}

func ConfigureLogger(opts *config.CliOptions) {
	verbose = opts.Verbose
}

func ErrorFmt(message string, vals... interface{}) {
	fmt.Fprintln(os.Stderr, nowStamp(), prefixError,
		nlPad(fmt.Sprintf(message, vals...)))
}

func Error(vals... interface{}) {
	fmt.Fprintln(os.Stderr, nowStamp(), prefixError, nlPad(fmt.Sprint(vals...)))
}

func WarnFmt(message string, vals... interface{}) {
	fmt.Fprintln(os.Stderr, nowStamp(), prefixWarn,
		nlPad(fmt.Sprintf(message, vals...)))
}

func Warn(vals... interface{}) {
	fmt.Fprintln(os.Stderr, nowStamp(), prefixWarn, nlPad(fmt.Sprint(vals...)))
}

func InfoFmt(message string, vals... interface{}) {
	fmt.Println(nowStamp(), prefixInfo, nlPad(fmt.Sprintf(message, vals...)))
}

func Info(vals... interface{}) {
	fmt.Println(nowStamp(), prefixInfo, nlPad(fmt.Sprint(vals...)))
}

func TraceFmt(message string, vals... interface{}) {
	if verbose {
		fmt.Println(nowStamp(), prefixDebug, nlPad(fmt.Sprintf(message, vals...)))
	}
}

func Trace(vals... interface{}) {
	if verbose {
		fmt.Println(nowStamp(), prefixDebug, nlPad(fmt.Sprint(vals...)))
	}
}
