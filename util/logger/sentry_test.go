package logger

import (
	"bytes"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSentryNoLevel(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	Println("abcdefg")
	require.True(t, strings.HasSuffix(out.String(), "abcdefg\n"), "Miss match value: %s", out.String())
	require.False(t, strings.Contains(out.String(), "[INFO]"), "Miss match value: %s", out.String())
	out.Reset()
}

func TestSentryPanic(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	defer func() {
		if err := recover(); err != nil {
			require.True(t, strings.HasSuffix(out.String(), "nonononon\n"), "Miss match value: %s", out.String())
			require.True(t, strings.Contains(out.String(), "[PANIC]"), "Miss match value: %s", out.String())
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

	require.True(t, strings.HasSuffix(out.String(), "ldkdkdkdks\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[CRITICAL]"), "Miss match value: %s", out.String())
	out.Reset()
}

func TestSentryError(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	E("lkerja;we")

	require.True(t, strings.HasSuffix(out.String(), "lkerja;we\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[ERROR]"), "Miss match value: %s", out.String())
	out.Reset()
}

func TestSentryWarn(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	W("jrlkaefj")

	require.True(t, strings.HasSuffix(out.String(), "jrlkaefj\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[WARN]"), "Miss match value: %s", out.String())
	out.Reset()
}

func TestSentryInfo(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	Infoln("abcdefg")
	require.True(t, strings.HasSuffix(out.String(), "abcdefg\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[INFO]"), "Miss match value: %s", out.String())
	out.Reset()
}

func TestSentryDebug(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	isDebug = true
	D("040itaokwp")

	require.True(t, strings.HasSuffix(out.String(), "040itaokwp\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[DEBUG]"), "Miss match value: %s", out.String())
	out.Reset()
	isDebug = false
}

func TestSentryTODO(t *testing.T) {
	t.Helper()

	out := &bytes.Buffer{}
	setup(out)

	T("lk2j3wr")

	require.True(t, strings.HasSuffix(out.String(), "lk2j3wr\n"), "Miss match value: %s", out.String())
	require.True(t, strings.Contains(out.String(), "[TODO]"), "Miss match value: %s", out.String())
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
	m.Run()
}
