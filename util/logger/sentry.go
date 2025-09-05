// Package logger ...
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

var (
	noLogger       = log.New(os.Stdout, "[NOLEVEL] ", log.LstdFlags|log.Llongfile)
	panicLogger    = log.New(os.Stdout, "[PANIC] ", log.LstdFlags|log.Llongfile)
	criticalLogger = log.New(os.Stdout, "[CRITICAL] ", log.LstdFlags|log.Llongfile)
	errLogger      = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Llongfile)
	warnLogger     = log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Llongfile)
	infoLogger     = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Llongfile)
	debugLogger    = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Llongfile)
	todoLogger     = log.New(os.Stdout, "[TODO] ", log.LstdFlags|log.Llongfile)
	isDebug        = false
	isSentry       = false
	isVerbose      = false
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

	// NOLEVEL Just for rename in sentry dashboard eventlog title
	NOLEVEL struct{ s string }
	// PANIC Just for rename in sentry dashboard eventlog title
	PANIC struct{ s string }
	// CRITICAL Just for rename in sentry dashboard eventlog title
	CRITICAL struct{ s string }
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
func (e *CRITICAL) Error() string { return e.s }
func (e *ERROR) Error() string    { return e.s }
func (e *WARN) Error() string     { return e.s }
func (e *INFO) Error() string     { return e.s }
func (e *DEBUG) Error() string    { return e.s }
func (e *TODO) Error() string     { return e.s }

// SetDebug set a debug by.
func SetDebug(debug bool) {
	isDebug = debug
}

// SetSentry a switcher
func SetSentry(sentry bool) {
	isSentry = sentry
}

// SetVerbose set a verbose logging by.
func SetVerbose(verbose bool) {
	isVerbose = verbose
}

// getSender returns the appropriate MessageFunc based on the context
func getSender(ctx context.Context) MessageFunc {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub.CaptureMessage
	}
	return sentry.CaptureMessage
}

// TodofWithContext outputs ...
func TodofWithContext(ctx context.Context, format string, args ...interface{}) {
	tododeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// TodolnWithContext outputs ...
func TodolnWithContext(ctx context.Context, args ...interface{}) {
	tododeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Todof outputs ...
func Todof(format string, args ...interface{}) {
	tododeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Todoln outputs ...
func Todoln(args ...interface{}) {
	tododeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	T = Todof // alias for Todof
	TCtx = TodofWithContext // alias for TodofWithContext
)

func tododeps(fn MessageFunc, deps int, s string) {
	_ = todoLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[TODO] %s \non %s:%d", s, f, line)
		// raven.CaptureMessage(s, nil, raven.NewException(&TODO{s}, trace(deps)))
		fn(s)
	}
}

// DebugfWithContext outputs ...
func DebugfWithContext(ctx context.Context, format string, args ...interface{}) {
	debugdeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// DebuglnWithContext outputs ...
func DebuglnWithContext(ctx context.Context, args ...interface{}) {
		debugdeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Debugf outputs ...
func Debugf(format string, args ...interface{}) {
	debugdeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Debugln outputs ...
func Debugln(args ...interface{}) {
	debugdeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	DCtx = DebugfWithContext // alias for DebugfWithContext
	D    = Debugf            // alias for Debugf
)

func debugdeps(fn MessageFunc, deps int, s string) {
	if isDebug {
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

// InfofWithContext pritns as information
func InfofWithContext(ctx context.Context, format string, args ...interface{}) {
	infodeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// InfolnWithContext pritns as information
func InfolnWithContext(ctx context.Context, args ...interface{}) {
	infodeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Infof pritns as information
func Infof(format string, args ...interface{}) {
	infodeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Infoln pritns as information
func Infoln(args ...interface{}) {
	infodeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	ICtx = InfofWithContext // alias for InfofWithContext
	I    = Infof            // alias for Infof
)


func infodeps(fn MessageFunc, deps int, s string) {
	_ = infoLogger.Output(deps, s)

	if isVerbose && isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[INFO] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&INFO{s}, trace(deps)))
		fn(s)
	}
}

// WarnfWithContext pritns as warning
func WarnfWithContext(ctx context.Context, format string, args ...interface{}) {
	warndeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// WarningfWithContext pritns as warning
func WarningfWithContext(ctx context.Context, format string, args ...interface{}) {
	warndeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// WarninglnWithContext outputs ...
func WarninglnWithContext(ctx context.Context, args ...interface{}) {
	warndeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// WarnlnWithContext outputs ...
func WarnlnWithContext(ctx context.Context, args ...interface{}) {
	warndeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Warningf pritns as warning
func Warningf(format string, args ...interface{}) {
	warndeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Warningln outputs ...
func Warningln(args ...interface{}) {
	warndeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	Warnf  = Warningf            // alias for Warningf
	Warnln = Warningln           // alias for Warningf
	WCtx   = WarningfWithContext // alias for WarninglnWithContext
	W      = Warningf            // alias for Warningf
)

func warndeps(fn MessageFunc, deps int, s string) {
	_ = warnLogger.Output(deps, s)

	if isVerbose && isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[WARN] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&WARN{s}, trace(deps)))
		fn(s)
	}
}

// ErrorlnWithContext prints as error
func ErrorlnWithContext(ctx context.Context, args ...interface{}) {
	errdeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// ErrorfWithContext prints as error
func ErrorfWithContext(ctx context.Context, format string, args ...interface{}) {
	errdeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// Errorln pritns as error
func Errorln(args ...interface{}) {
	errdeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

// Errorf pritns as error
func Errorf(format string, args ...interface{}) {
	errdeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

var (
	ErrWithContext   = ErrorfWithContext  // alias for ErrorfWithContext
	ErrlnWithContext = ErrorlnWithContext // alias for ErrorlnWithContext
	Errln            = Errorln            // alias for Errorln
	Errf             = Errorf             // alias for Errorf
	ECtx             = ErrorfWithContext  // alias for ErrorfWithContext
	E                = Errorf             // alias for Errorf
)

func errdeps(fn MessageFunc, deps int, s string) {
	_ = errLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[ERROR] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&ERROR{s}, trace(deps)))
		fn(s)
	}
}

// CriticalfWithContext prints as critical
func CriticalfWithContext(ctx context.Context, format string, args ...interface{}) {
	criticaldeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// CriticalnWithContext outputs ...
func CriticalnWithContext(ctx context.Context, args ...interface{}) {
	criticaldeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Criticalf pritns as critical
func Criticalf(format string, args ...interface{}) {
	criticaldeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Criticaln outputs ...
func Criticaln(args ...interface{}) {
	criticaldeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	CrtlfWithContext = CriticalfWithContext // alias for CriticalfWithContext
	CCtx             = CriticalfWithContext // alias for CriticalfWithContext
	CrtlnWithContext = CriticalnWithContext // alias for CriticalnWithContext
	Crtln            = Criticaln            // alias for Criticaln
	Crtlf            = Criticalf            // alias for Criticalf
	C                = Criticalf            // alias for Criticalf
)

func criticaldeps(fn MessageFunc, deps int, s string) {
	_ = criticalLogger.Output(deps, s)

	if isSentry {
		_, f, line, _ := runtime.Caller(deps - 1)
		s = fmt.Sprintf("[CRITICAL] %s \non %s:%d", s, f, line)
		// TODO: tidy up
		// raven.CaptureMessage(s, nil, sentry.NewException(&CRITICAL{s}, trace(deps)))
		fn(s)
	}
}

// PanicfWithContext prints as panic
func PanicfWithContext(ctx context.Context, format string, args ...interface{}) {
	panicdeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// PaniclnWithContext outputs ...
func PaniclnWithContext(ctx context.Context, args ...interface{}) {
	panicdeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Panicf pritns as panic
func Panicf(format string, args ...interface{}) {
	panicdeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Panicln outputs ...
func Panicln(args ...interface{}) {
	panicdeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

var (
	PCtx = PanicfWithContext // alias for PanicfWithContext
	P    = Panicf            // alias for Panicf
)

func panicdeps(fn MessageFunc, deps int, s string) {
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

// PrintfWithContext prints with format
func PrintfWithContext(ctx context.Context, format string, args ...interface{}) {
	printdeps(getSender(ctx), 3, fmt.Sprintf(format, args...))
}

// PrintlnWithContext outputs ...
func PrintlnWithContext(ctx context.Context, args ...interface{}) {
	printdeps(getSender(ctx), 3, fmt.Sprintln(args...))
}

// Printf pritns with format
func Printf(format string, args ...interface{}) {
	printdeps(sentry.CaptureMessage, 3, fmt.Sprintf(format, args...))
}

// Println outputs ...
func Println(args ...interface{}) {
	printdeps(sentry.CaptureMessage, 3, fmt.Sprintln(args...))
}

func printdeps(fn MessageFunc, deps int, s string) {
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
