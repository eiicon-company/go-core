package dsn

import (
	"net/url"
	"strings"
)

// RedisDSN redis://127.0.0.1:6379/8
type RedisDSN struct {
	Password, Host, Port, HostPort, DB string
}

// Redis ...
func Redis(uri string) (*RedisDSN, error) {
	if uri == "" {
		return nil, ef("invalid redis dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, ef("invalid redis hasn't scheme")
	}

	hp := strings.Split(u.Host, ":")
	if len(hp) != 2 {
		return nil, ef("invalid redis host:port: %s", u.Host)
	}
	db := strings.Split(u.Path, "/")
	if len(db) != 2 {
		return nil, ef("invalid redis db number: %s", u.Path)
	}

	dsn := &RedisDSN{}

	if u.User != nil {
		password, ok := u.User.Password()
		if ok {
			dsn.Password = password
		}
	}

	dsn.DB = db[1]
	dsn.HostPort = u.Host
	dsn.Host, dsn.Port = hp[0], hp[1]

	return dsn, nil
}

// Flags ...
func (c *RedisDSN) Flags() []string {
	flags := []string{}
	if c.Password != `` {
		flags = append(flags, `-a`, c.Password)
	}
	flags = append(flags, `-h`, c.Host, `-p`, c.Port, `-n`, c.DB)
	return flags
}
