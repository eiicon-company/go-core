package testdb

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"

	"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql/driver"
)

type (
	mysqlTester struct {
		name    string
		host    string
		user    string
		pass    string
		port    int
		sslmode bool
		schema  []byte // Table DDL

		db *sql.DB
	}
)

func (m *mysqlTester) cmdArgs() []string {
	return []string{
		fmt.Sprintf("-u%s", m.user),
		fmt.Sprintf("--password=%s", m.pass),
		fmt.Sprintf("-h%s", m.host),
		fmt.Sprintf("-P%d", m.port),
	}
}

func (m *mysqlTester) createDB() error {
	sql := fmt.Sprintf("CREATE DATABASE %s;", m.name)
	return m.cmd(sql)
}

func (m *mysqlTester) dropDB() error {
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", m.name)
	return m.cmd(sql)
}

func (m *mysqlTester) Setup() error {
	if err := m.dropDB(); err != nil {
		return err
	}
	if err := m.createDB(); err != nil {
		return err
	}

	if err := m.cmd(string(m.schema), "-D", m.name); err != nil {
		return xerrors.Errorf("failed table restore: %w", err)
	}

	return nil
}

func (m *mysqlTester) Teardown() error {
	if m.db != nil {
		if err := m.db.Close(); err != nil {
			return err
		}
	}

	if err := m.dropDB(); err != nil {
		return err
	}

	return nil
}

// StdinCommand execs database
func StdinCommand(stdin string, args ...string) error {
	cmd := exec.Command(dialect, args...)
	cmd.Stdin = strings.NewReader(stdin)

	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr

	if err := cmd.Run(); err != nil {
		return xerrors.Errorf("failed running cmd=%s args=%s out=%s err=%s: %w",
			stdin, args, stdout.String(), stderr.String(), err)
	}

	_ = cmd.Wait() //#nosec no need error

	return nil
}

func (m *mysqlTester) cmd(stdin string, args ...string) error {
	args = append(m.cmdArgs(), args...)
	return StdinCommand(stdin, args...)
}

func (m *mysqlTester) Conn() (*sql.DB, error) {
	if m.db != nil {
		return m.db, nil
	}

	dsn := driver.MySQLBuildQueryString(
		m.user, m.pass, m.name, m.host, m.port, fmt.Sprint(m.sslmode))

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	m.db = db
	return m.db, nil
}
