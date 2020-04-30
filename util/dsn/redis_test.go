package dsn

import (
	"testing"
)

func TestRedis(t *testing.T) {
	t.Helper()

	f, err := Redis("file://127.0.0.1:6379/4")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = Redis("redis://127.0.0.1:6379/4")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	if f.Password != "" {
		t.Error("redis field error")
	}
	if f.Host != "127.0.0.1" {
		t.Error("redis field error")
	}
	if f.Port != "6379" {
		t.Error("redis field error")
	}
	if f.HostPort != "127.0.0.1:6379" {
		t.Error("redis field error")
	}
	if f.DB != "4" {
		t.Error("redis field error")
	}

	t.Logf("Redis: %#+v", f)
}
