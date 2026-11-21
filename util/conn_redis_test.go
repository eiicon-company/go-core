package util

import (
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/mediocregopher/radix/v3"
)

// redisEnv answers EnvString for the Redis keys the connectors read; the
// embedded Environment is never called.
type redisEnv struct {
	Environment
	values map[string]string
}

func (e redisEnv) EnvString(key string) string { return e.values[key] }

func dlmEnv(uri string) redisEnv { return redisEnv{values: map[string]string{"DLMURI": uri}} }

// TestSelectRedisConnSendsPassword guards that a password in the URI is sent as
// AUTH, and that a server requiring one refuses a URI without it.
func TestSelectRedisConnSendsPassword(t *testing.T) {
	// Arrange
	server := miniredis.RunT(t)
	server.RequireAuth("s3cret")

	// Act
	pool, err := SelectRedisConn("redis://:s3cret@" + server.Addr() + "/3")

	// Assert
	if err != nil {
		t.Fatalf("Unable to connect with the password: %s", err)
	}
	defer pool.Close()
	var pong string
	if err := pool.Do(radix.Cmd(&pong, "PING")); err != nil || pong != "PONG" {
		t.Fatalf("Miss match value: want PONG, got %q (%v)", pong, err)
	}

	if _, err := SelectRedisConn("redis://" + server.Addr() + "/3"); err == nil {
		t.Fatal("Connected without the password the server requires")
	}
}

// TestDLMConnSendsPassword is the same guard for the distributed lock pool.
func TestDLMConnSendsPassword(t *testing.T) {
	// Arrange
	server := miniredis.RunT(t)
	server.RequireAuth("s3cret")

	// Act
	_, err := DLMConn(dlmEnv("redis://:s3cret@" + server.Addr() + "/2"))

	// Assert
	if err != nil {
		t.Fatalf("Unable to connect with the password: %s", err)
	}
	if _, err := DLMConn(dlmEnv("redis://" + server.Addr() + "/2")); err == nil {
		t.Fatal("Connected without the password the server requires")
	}
}

// TestRedisConnWithoutPasswordSendsNoAuth guards the servers that have no
// password: sending AUTH to them would fail.
func TestRedisConnWithoutPasswordSendsNoAuth(t *testing.T) {
	// Arrange
	server := miniredis.RunT(t)

	// Act
	pool, err := SelectRedisConn("redis://" + server.Addr() + "/1")
	_, dlmErr := DLMConn(dlmEnv("redis://" + server.Addr() + "/1"))

	// Assert
	if err != nil || dlmErr != nil {
		t.Fatalf("Unable to connect without a password: %v, %v", err, dlmErr)
	}
	pool.Close()
}

// TestRedisConnParseErrorHidesPassword guards the error of a URI that does
// not parse: url.Parse repeats the whole URI in its own error.
func TestRedisConnParseErrorHidesPassword(t *testing.T) {
	const uri = "redis://:s3cret@bad host:6379/8"

	_, err := SelectRedisConn(uri)
	_, dlmErr := DLMConn(dlmEnv(uri))

	for _, e := range []error{err, dlmErr} {
		if e == nil {
			t.Fatal("Parsed a URI with a space in its host")
		}
		if strings.Contains(e.Error(), "s3cret") {
			t.Errorf("Password left in the error: %s", e)
		}
	}
}

// TestRedactRedisURIHidesPassword guards the connection URIs written to logs
// and errors.
func TestRedactRedisURIHidesPassword(t *testing.T) {
	cases := map[string]string{ //nolint:gosec // G101: fake passwords the test must find masked
		"redis://:s3cret@127.0.0.1:6379/8":     "redis://:xxxxx@127.0.0.1:6379/8",
		"redis://user:s3cret@127.0.0.1:6379/8": "redis://user:xxxxx@127.0.0.1:6379/8",
		"redis://127.0.0.1:6379/8":             "redis://127.0.0.1:6379/8",
		"redis://:s3cret@bad host:6379/8":      "<unparsable redis uri>",
	}
	for uri, want := range cases {
		got := redactRedisURI(uri)
		if got != want {
			t.Errorf("Miss match value: redactRedisURI(%q) = %q, want %q", uri, got, want)
		}
		if strings.Contains(got, "s3cret") {
			t.Errorf("Password left in %q", got)
		}
	}
}
