package dsn

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

// MailDSN stdout:// or smtp://username@gmail.com:password@smtp.gmail.com(smtp.gmail.com:587)/?tls=false
type MailDSN struct {
	// Auth
	User, Password, Host string
	// Server
	Addr string
	// Option
	TLS, StdOut bool
}

// Mail stdout:// or smtp://username@gmail.com:password@smtp.gmail.com(smtp.gmail.com:587)/?tls=false
func Mail(uri string) (*MailDSN, error) {
	if strings.HasPrefix(uri, "stdout://") {
		return &MailDSN{StdOut: true}, nil
	}

	if uri == "" {
		return nil, ef("invalid mail dsn")
	}
	if !strings.HasPrefix(uri, "smtp://") {
		return nil, ef("invalid mail scheme. e.g. smtp://")
	}

	m, err := mysql.ParseDSN(strings.TrimPrefix(uri, "smtp://"))
	if err != nil {
		return nil, err
	}
	if m.User == "" {
		return nil, ef("invalid mail hasn't auth user: %s", m.User)
	}
	if m.Passwd == "" {
		return nil, ef("invalid mail hasn't auth password")
	}
	if m.Net == "" {
		return nil, ef("invalid mail hasn't auth host: %s", m.Net)
	}
	hp := strings.Split(m.Addr, ":")
	if len(hp) != 2 {
		return nil, ef("invalid mail host:port: %s", m.Addr)
	}

	dsn := &MailDSN{
		User:     m.User,
		Password: m.Passwd,
		Host:     m.Net,
		Addr:     m.Addr,
		TLS:      m.TLSConfig == "true",
	}

	return dsn, nil
}
