// Package logger is simply logger with  sentry
package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

var (
	noLogger       = log.New(os.Stdout, "[NOLEVEL] ", log.LstdFlags|log.Llongfile)
	panicLogger    = log.New(os.Stdout, "[PANIC] ", log.LstdFlags|log.Llongfile)
	creticalLogger = log.New(os.Stdout, "[CRETICAL] ", log.LstdFlags|log.Llongfile)
	errLogger      = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Llongfile)
	warnLogger     = log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Llongfile)
	infoLogger     = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Llongfile)
	debugLogger    = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Llongfile)
	todoLogger     = log.New(os.Stdout, "[TODO] ", log.LstdFlags|log.Llongfile)
	isDebug        = false
	isSentry       = false
	without        = false
)

// var (
// 	major = []string{"10", "11", "12", "13", "14", "15", "16", "17"}
// 	minor = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
// 	vers  = []string{"github.com"}
// )

func init() {
	// SetAttachStacktrace(true)

	// for _, ma := range major {
	// 	for _, mi := range minor {
	// 		vers = append(vers, fmt.Sprintf("/root/.gvm/pkgsets/go1.%s.%s/global/src/github.com/eiicon-company", ma, mi))
	// 		vers = append(vers, fmt.Sprintf("/root/.gvm/pkgsets/go1.%s.%s/global/pkg/mod/github.com/eiicon-company", ma, mi))
	// 		vers = append(vers, fmt.Sprintf("%s/src/github.com/eiicon-company", os.Getenv("GOPATH")))
	// 		vers = append(vers, fmt.Sprintf("%s/pkg/mod/github.com/eiicon-company", os.Getenv("GOPATH")))
	// 	}
	// }
}

type (
	// NOLEVEL Just for rename in sentry dashboard eventlog title
	NOLEVEL struct{ s string }
	// PANIC Just for rename in sentry dashboard eventlog title
	PANIC struct{ s string }
	// CRETICAL Just for rename in sentry dashboard eventlog title
	CRETICAL struct{ s string }
	// ERROR Just for rename in sentry dashboard eventlog title
	ERROR struct{ s string }
	// WARN Just for rename in sentry dashboard eventlog title
	WARN struct{ s string }
	// INFO Just for rename in sentry dashboard eventlog title
	INFO struct{ s string }
	// DEBUG Just for rename in sentry dashboard eventlog title
	DEBUG struct{ s string }
	// TODO Just for rename in sentry dashboard eventlog title
	TODO struct{ s string }
)

func (e *NOLEVEL) Error() string  { return e.s }
func (e *PANIC) Error() string    { return e.s }
func (e *CRETICAL) Error() string { return e.s }
func (e *ERROR) Error() string    { return e.s }
func (e *WARN) Error() string     { return e.s }
func (e *INFO) Error() string     { return e.s }
func (e *DEBUG) Error() string    { return e.s }
func (e *TODO) Error() string     { return e.s }

// C is alias with cretical
func C(format string, args ...interface{}) { creticaldeps(3, format, args...) }

// E is alias with error
func E(format string, args ...interface{}) { errdeps(3, format, args...) }

// W is alias with warning
func W(format string, args ...interface{}) { warndeps(3, format, args...) }

// I is alias with info
func I(format string, args ...interface{}) { infodeps(3, format, args...) }

// D is alias with debug
func D(format string, args ...interface{}) { debugdeps(3, format, args...) }

// T is alias with todo
func T(format string, args ...interface{}) { tododeps(3, format, args...) }

// SetDebug set debug flag
func SetDebug(debug bool) {
	isDebug = debug
}

// SetSentry set
func SetSentry(sentry bool) {
	isSentry = sentry
}

// func trace(deps int) *sentry.Stacktrace {
// 	return sentry.NewStacktrace() // TODO: filter deps, 5, vers
// }

// Todo outputs ...
func Todo(format string, args ...interface{}) {
	tododeps(3, format, args...)
}

func tododeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = todoLogger.Output(deps, s)

	if isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[TODO] %s:%d: %s", fn, line, s)

		// scope.SetExtras(map[string]interface{}{"path": path, "cwd": os.Getwd()})
		// nil, sentry.NewException(&TODO{s}, trace(deps))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}
}

// Debug outputs ...
func Debug(format string, args ...interface{}) {
	debugdeps(3, format, args...)
}

func debugdeps(deps int, format string, args ...interface{}) {
	if isDebug {
		s := fmt.Sprintf(format, args...)
		_ = debugLogger.Output(deps, s)

		if isSentry {
			_, fn, line, _ := runtime.Caller(deps - 1)
			s = fmt.Sprintf("[DEBUG] %s:%d: %s", fn, line, s)
			// sentry.CaptureMessage(s, nil, sentry.NewException(&DEBUG{s}, trace(deps)))
			sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
		}
	}
}

// Infof pritns as information
func Infof(format string, args ...interface{}) {
	infodeps(3, format, args...)
}

// Info pritns as information
func Info(format string, args ...interface{}) {
	infodeps(3, format, args...)
}

func infodeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = infoLogger.Output(deps, s)

	if without && isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[INFO] %s:%d: %s", fn, line, s)
		// sentry.CaptureMessage(s, nil, sentry.NewException(&INFO{s}, trace(deps)))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}
}

// Warnf pritns as warning
func Warnf(format string, args ...interface{}) {
	warndeps(3, format, args...)
}

// Warningf pritns as warning
func Warningf(format string, args ...interface{}) {
	warndeps(3, format, args...)
}

// Warn outputs ...
func Warn(format string, args ...interface{}) {
	warndeps(3, format, args...)
}

func warndeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = warnLogger.Output(deps, s)

	if without && isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[WARN] %s:%d: %s", fn, line, s)
		// sentry.CaptureMessage(s, nil, sentry.NewException(&WARN{s}, trace(deps)))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}
}

// Error pritns as error
func Error(format string, args ...interface{}) {
	errdeps(3, format, args...)
}

// Errorf pritns as error
func Errorf(format string, args ...interface{}) {
	errdeps(3, format, args...)
}

// Errf pritns as error
func Errf(format string, args ...interface{}) {
	errdeps(3, format, args...)
}

// Err outputs ...
func Err(format string, args ...interface{}) {
	errdeps(3, format, args...)
}

func errdeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = errLogger.Output(deps, s)

	if isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[ERROR] %s:%d: %s", fn, line, s)
		// sentry.CaptureMessage(s, nil, sentry.NewException(&ERROR{s}, trace(deps)))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}
}

// Creticalf pritns as cretical
func Creticalf(format string, args ...interface{}) {
	creticaldeps(3, format, args...)
}

// Cretical outputs ...
func Cretical(format string, args ...interface{}) {
	creticaldeps(3, format, args...)
}

// Crtl pritns as cretical
func Crtl(format string, args ...interface{}) {
	creticaldeps(3, format, args...)
}

func creticaldeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = creticalLogger.Output(deps, s)

	if isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[CRETICAL] %s:%d: %s", fn, line, s)
		// sentry.CaptureMessage(s, nil, sentry.NewException(&CRETICAL{s}, trace(deps)))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}
}

// Panicf pritns as panic
func Panicf(format string, args ...interface{}) {
	panicdeps(3, format, args...)
}

// Panic outputs ...
func Panic(format string, args ...interface{}) {
	panicdeps(3, format, args...)
}

func panicdeps(deps int, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_ = panicLogger.Output(deps, s)

	if isSentry {
		_, fn, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[PANIC] %s:%d: %s", fn, line, s)
		// sentry.CaptureMessage(s, nil, sentry.NewException(&PANIC{s}, trace(deps)))
		sentry.WithScope(func(scope *sentry.Scope) { sentry.CaptureMessage(s) })
	}

	panic(s)
}

// Printf pritns with format
func Printf(format string, args ...interface{}) {
	printdeps(3, fmt.Sprintf(format, args...))
}

// Println outputs ...
func Println(args ...interface{}) {
	printdeps(3, args...)
}

func printdeps(deps int, args ...interface{}) {
	s := fmt.Sprintln(args...)

	switch {
	case strings.Contains(s, "[PANIC]"):
		creticaldeps(deps+1, s) // TODO: soft panic logging
	case strings.Contains(s, "[CRETICAL]"):
		creticaldeps(deps+1, s)
	case strings.Contains(s, "[ERROR]"):
		errdeps(deps+1, s)
	case strings.Contains(s, "[WARN]"):
		warndeps(deps+1, s)
	case strings.Contains(s, "[INFO]"):
		infodeps(deps+1, s)
	case strings.Contains(s, "[DEBUG]"):
		debugdeps(deps+1, s)
	case strings.Contains(s, "[TODO]"):
		tododeps(deps+1, s)
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
