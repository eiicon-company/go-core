package dsn

import (
	"testing"
)

func TestMailStdOut(t *testing.T) {
	f, err := Mail("stdout://")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	if !f.StdOut {
		t.Fatalf("Unknown Scheme: stdout is missing")
	}
}

func TestMail(t *testing.T) {
	t.Helper()

	f, err := Mail("smtp://username@gmail.com:password@smtp.gmail.com(smtp.gmail.com:587)/?tls=false")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	if f.User != "username@gmail.com" {
		t.Error("mail field error")
	}
	if f.Password != "password" {
		t.Error("mail field error")
	}
	if f.Host != "smtp.gmail.com" {
		t.Error("mail field error")
	}
	if f.Addr != "smtp.gmail.com:587" {
		t.Error("mail field error")
	}
	if f.TLSServer != "smtp.gmail.com" {
		t.Error("mail field error")
	}
	if f.TLS != false {
		t.Error("mail field error")
	}
	if f.StdOut != false {
		t.Error("mail field error")
	}

	t.Logf("Mail: %#+v", f)
}
