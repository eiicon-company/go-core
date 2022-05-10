package logger

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSentryNoLevel(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	Println("abcdefg")
	if !strings.HasSuffix(out.String(), "abcdefg\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if strings.Contains(out.String(), "[INFO]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func TestSentryPanic(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	defer func() {
		if err := recover(); err != nil {
			if !strings.HasSuffix(out.String(), "nonononon\n") {
				t.Errorf("Miss match value: %s", out.String())
			}
			if !strings.Contains(out.String(), "[PANIC]") {
				t.Errorf("Miss match value: %s", out.String())
			}
			out.Reset()
		}
	}()

	Panicln("nonononon")
}

func TestSentryCritical(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	C("ldkdkdkdks")

	if !strings.HasSuffix(out.String(), "ldkdkdkdks\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[CRITICAL]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func TestSentryError(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	E("lkerja;we")

	if !strings.HasSuffix(out.String(), "lkerja;we\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[ERROR]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func TestSentryWarn(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	W("jrlkaefj")

	if !strings.HasSuffix(out.String(), "jrlkaefj\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[WARN]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func TestSentryInfo(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	Infoln("abcdefg")
	if !strings.HasSuffix(out.String(), "abcdefg\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[INFO]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func TestSentryDebug(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	isDebug = true
	D("040itaokwp")

	if !strings.HasSuffix(out.String(), "040itaokwp\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[DEBUG]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
	isDebug = false
}

func TestSentryTODO(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	T("lk2j3wr")

	if !strings.HasSuffix(out.String(), "lk2j3wr\n") {
		t.Errorf("Miss match value: %s", out.String())
	}
	if !strings.Contains(out.String(), "[TODO]") {
		t.Errorf("Miss match value: %s", out.String())
	}
	out.Reset()
}

func setup(out io.Writer) {
	noLogger = log.New(out, "[NOLEVEL] ", log.LstdFlags|log.Llongfile)
	panicLogger = log.New(out, "[PANIC] ", log.LstdFlags|log.Llongfile)
	criticalLogger = log.New(out, "[CRITICAL] ", log.LstdFlags|log.Llongfile)
	errLogger = log.New(out, "[ERROR] ", log.LstdFlags|log.Llongfile)
	warnLogger = log.New(out, "[WARN] ", log.LstdFlags|log.Llongfile)
	infoLogger = log.New(out, "[INFO] ", log.LstdFlags|log.Llongfile)
	debugLogger = log.New(out, "[DEBUG] ", log.LstdFlags|log.Llongfile)
	todoLogger = log.New(out, "[TODO] ", log.LstdFlags|log.Llongfile)
}

func TestMain(m *testing.M) {
	origNoLogger := noLogger
	origPanicLogger := panicLogger
	origCriticalLogger := criticalLogger
	origErrLogger := errLogger
	origWarnLogger := warnLogger
	origInfoLogger := infoLogger
	origDebugLogger := debugLogger
	origTodoLogger := todoLogger

	code := m.Run()

	noLogger = origNoLogger
	panicLogger = origPanicLogger
	criticalLogger = origCriticalLogger
	errLogger = origErrLogger
	warnLogger = origWarnLogger
	infoLogger = origInfoLogger
	debugLogger = origDebugLogger
	todoLogger = origTodoLogger

	os.Exit(code)
}
