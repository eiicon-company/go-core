// Package logger is simply logger with sentry
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

const defaultDepth = 3

var (
	noLogger              = log.New(os.Stdout, "[NOLEVEL] ", log.LstdFlags|log.Llongfile)
	panicLogger           = log.New(os.Stdout, "[PANIC] ", log.LstdFlags|log.Llongfile)
	criticalLogger        = log.New(os.Stdout, "[CRITICAL] ", log.LstdFlags|log.Llongfile)
	errLogger             = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Llongfile)
	warnLogger            = log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Llongfile)
	infoLogger            = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Llongfile)
	debugLogger           = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Llongfile)
	todoLogger            = log.New(os.Stdout, "[TODO] ", log.LstdFlags|log.Llongfile)
	isDebug               = false
	isSentry              = false
	isVerbose             = false
	defaultCaptureMessage = sentry.CaptureMessage
)

// TODO: tidy up
// var (
// 	major = []string{"10", "11", "12", "13", "14", "15", "16", "17"}
// 	minor = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
// 	vers  = []string{"github.com"}
// )

// func init() {
// 	// TODO: tidy up
// 	SetAttachStacktrace(true)
//
// 	for _, ma := range major {
// 		for _, mi := range minor {
// 			vers = append(vers, fmt.Sprintf("/root/.gvm/pkgsets/go1.%s.%s/global/src/github.com/eiicon-company", ma, mi))
// 			vers = append(vers, fmt.Sprintf("/root/.gvm/pkgsets/go1.%s.%s/global/pkg/mod/github.com/eiicon-company", ma, mi))
// 			vers = append(vers, fmt.Sprintf("%s/src/github.com/eiicon-company", os.Getenv("GOPATH")))
// 			vers = append(vers, fmt.Sprintf("%s/pkg/mod/github.com/eiicon-company", os.Getenv("GOPATH")))
// 		}
// 	}
// }

// TODO: tidy up
// func trace(deps int) *sentry.Stacktrace {
// 	return sentry.NewStacktrace() // TODO: filter deps, 5, vers
// }

type (
	// MessageFunc uses both func & method(hub.XXXX, sentry.XXXX)
	MessageFunc func(message string) *sentry.EventID
)

// getMessageFuncFromContext fetches MessageFunc from context, or defaultCaptureMessage.
func getMessageFuncFromContext(ctx context.Context) MessageFunc {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureMessage
	}
	return defaultCaptureMessage
}

// CriticalfWithContext pritns as critical
func CriticalfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	criticaldeps(msg, defaultDepth, format, args...)
}

// ErrorfWithContext pritns as error
func ErrorfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	errdeps(msg, defaultDepth, format, args...)
}

// WarnfWithContext pritns as warn
func WarnfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	warndeps(msg, defaultDepth, format, args...)
}

// InfofWithContext pritns as information
func InfofWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	infodeps(msg, defaultDepth, format, args...)
}

// DebugfWithContext pritns as debug
func DebugfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	debugdeps(msg, defaultDepth, format, args...)
}

// TodofWithContext pritns as todo
func TodofWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	tododeps(msg, defaultDepth, format, args...)
}

// PanicfWithContext prints as panic
func PanicfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	panicdeps(msg, defaultDepth, format, args...)
}

// PrintfWithContext prints with format
func PrintfWithContext(ctx context.Context, format string, args ...interface{}) {
	msg := getMessageFuncFromContext(ctx)

	printdeps(msg, defaultDepth, format, args...)
}

// Todof outputs ...
func Todof(format string, args ...interface{}) {
	tododeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Debugf outputs ...
func Debugf(format string, args ...interface{}) {
	debugdeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Infof pritns as information
func Infof(format string, args ...interface{}) {
	infodeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Warnf pritns as warning
func Warnf(format string, args ...interface{}) {
	warndeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Errorf pritns as error
func Errorf(format string, args ...interface{}) {
	errdeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Criticalf pritns as critical
func Criticalf(format string, args ...interface{}) {
	criticaldeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Panicf pritns as panic
func Panicf(format string, args ...interface{}) {
	panicdeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// Printf pritns with format
func Printf(format string, args ...interface{}) {
	printdeps(defaultCaptureMessage, defaultDepth, format, args...)
}

// SetDebug set a debug by.
func SetDebug(debug bool) {
	isDebug = debug
}

// SetSentry a switcher
func SetSentry(sentry bool) {
	isSentry = sentry
}

// SetVerbose set a verbose logging by.
func SetVerbose(debug bool) {
	isVerbose = debug
}

func tododeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = todoLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[TODO] %s \non %s:%d", s, f, line)
		// raven.CaptureMessage(s, nil, raven.NewException(&TODO{s}, trace(deps)))
		fn(s)
	}
}

func debugdeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	if isDebug {
		s := fmt.Sprintf(format, args...)
		_ = debugLogger.Output(deps, s)

		if isVerbose && isSentry {
			_, f, line, _ := runtime.Caller(deps - 1)
			s = fmt.Sprintf("[DEBUG] %s \non %s:%d", s, f, line)
			// TODO: tidy up
			// raven.CaptureMessage(s, nil, sentry.NewException(&DEBUG{s}, trace(deps)))
			fn(s)
		}
	}
}

func infodeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = infoLogger.Output(deps, s)

	if isVerbose && isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[INFO] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&INFO{s}, trace(deps)))
		fn(s)
	}
}

func warndeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = warnLogger.Output(deps, s)

	if isVerbose && isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[WARN] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&WARN{s}, trace(deps)))
		fn(s)
	}
}

func errdeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = errLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[ERROR] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&ERROR{s}, trace(deps)))
		fn(s)
	}
}

func criticaldeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = criticalLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[CRITICAL] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&CRITICAL{s}, trace(deps)))
		fn(s)
	}
}

func panicdeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = panicLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[PANIC] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&PANIC{s}, trace(deps)))
		fn(s)
	}

	panic(s)
}

//nolint:govet
func printdeps(fn MessageFunc, deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)

	switch {
	case strings.Contains(s, "[PANIC]"):
		criticaldeps(fn, deps+1, s) // TODO: soft panic logging
	case strings.Contains(s, "[CRITICAL]"):
		criticaldeps(fn, deps+1, s)
	case strings.Contains(s, "[ERROR]"):
		errdeps(fn, deps+1, s)
	case strings.Contains(s, "[WARN]"):
		warndeps(fn, deps+1, s)
	case strings.Contains(s, "[INFO]"):
		infodeps(fn, deps+1, s)
	case strings.Contains(s, "[DEBUG]"):
		debugdeps(fn, deps+1, s)
	case strings.Contains(s, "[TODO]"):
		tododeps(fn, deps+1, s)
	default:
		_ = noLogger.Output(deps, s)
	}
}

// Flush waits blocks and waits for all events to finish being sent to Sentry server
func Flush() {
	sentry.Flush(3 * time.Second)
}

// FlushTimeout waits blocks and waits for all events to finish being sent to Sentry server
func FlushTimeout(timeout time.Duration) {
	sentry.Flush(timeout)
}
