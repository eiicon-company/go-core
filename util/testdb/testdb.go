package testdb

import (
	"database/sql"
	"github.com/volatiletech/randomize"
	"net"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/xerrors"
)

var (
	dialect string
)

type (
	// DBTester enforces implementation entire interface's method.
	DBTester interface {
		Setup() error
		Teardown() error
		Conn() (*sql.DB, error)
	}
)

// NewDBTester returns a DBTester
func NewDBTester(dsn string, schema []byte) (DBTester, error) {
	if strings.HasPrefix(dsn, "mysql") {
		dialect = "mysql"

		my, err := mysql.ParseDSN(strings.TrimPrefix(dsn, "mysql://"))
		if err != nil {
			return nil, xerrors.Errorf("%sTester parses error: %w", dialect, err)
		}

		host, sport, err := net.SplitHostPort(my.Addr)
		if err != nil {
			return nil, xerrors.Errorf("%sTester parses error: %w", dialect, err)
		}

		port, err := strconv.Atoi(sport)
		if err != nil {
			return nil, xerrors.Errorf("%sTester parses error: %w", dialect, err)
		}

		return &mysqlTester{
			name:    randomize.StableDBName(my.DBName),
			host:    host,
			user:    my.User,
			pass:    my.Passwd,
			port:    port,
			sslmode: my.TLSConfig != "",
			schema:  schema,
		}, nil
	}

	return nil, xerrors.New("dbtester has no dialect")
}
