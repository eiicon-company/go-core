package util

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/getsentry/sentry-go"
	"github.com/go-sql-driver/mysql"
	dlmredis "github.com/gomodule/redigo/redis"
	radix "github.com/mediocregopher/radix/v3"
	proxy "github.com/shogo82148/go-sql-proxy"
	"github.com/spf13/cast"

	"github.com/eiicon-company/go-core/util/dlm"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
)

// DBConn returns current database established connection
func DBConn(dialect string, env Environment) (*sql.DB, error) {
	return SelectDBConn(dialect, env.EnvString("DSN"))
}

// SelectDBConn can choose db connection
func SelectDBConn(dialect, dsn string) (*sql.DB, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, fmt.Errorf("it was unable to connect the DB. %s", err)
	}

	// db configuration
	//
	// https://github.blog/2020-05-20-three-bugs-in-the-go-mysql-driver/
	// Oh Gawd
	// - https://github.com/go-sql-driver/mysql/issues/1302#issuecomment-1019842712
	// - https://github.com/go-sql-driver/mysql/issues/1120#issuecomment-636795680
	//
	//
	// https://github.com/go-sql-driver/mysql?tab=readme-ov-file#important-settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(6)
	db.SetMaxOpenConns(6)

	// make sure connection available
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("it was unable to connect the DB: %s", err)
	}

	var ver string
	logger.D("%s", db.QueryRow("SELECT @@version").Scan(&ver))

	msg := "[INFO] the mysql connection established <%s>, version %s"
	logger.Printf(msg, strings.Join(strings.Split(dsn, "@")[1:], ""), ver)

	return db, nil
}

// DBSlowQuery applies it with sentry span
//
// https://github.com/getsentry/sentry-ruby/issues/1674
// https://develop.sentry.dev/sdk/performance/span-operations/#database
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/database.md
func DBSlowQuery(dialect string, period time.Duration) {
	sql.Register(dialect, proxy.NewProxyContext(&mysql.MySQLDriver{}, &proxy.HooksContext{
		PreExec: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostExec: func(ctx context.Context, dt interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, _ error) error {
			startTime := dt.(time.Time)
			emitSlowSpan(ctx, "db.sql.exec.slow", startTime, time.Since(startTime), period, stmt.QueryString, args)

			return nil
		},
		PreQuery: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostQuery: func(ctx context.Context, dt interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Rows, _ error) error {
			startTime := dt.(time.Time)
			emitSlowSpan(ctx, "db.sql.query.slow", startTime, time.Since(startTime), period, stmt.QueryString, args)

			return nil
		},
	}))
}

// emitSlowSpan reports a statement that took longer than period to sentry, as a
// span nested in ctx. It is a no-op for statements that finished in time.
//
// since must be the measured elapsed time of the statement: the span is dated
// [startTime, startTime+since), so passing a since that disagrees with startTime
// misreports the duration.
func emitSlowSpan(ctx context.Context, op string, startTime time.Time, since, period time.Duration, query string, args []driver.NamedValue) {
	if since <= period {
		return
	}

	span := sentry.StartSpan(ctx, op, func(s *sentry.Span) {
		s.StartTime = startTime
		s.EndTime = startTime.Add(since)
		s.Description = query

		data := map[string]interface{}{}
		for i, arg := range args {
			if i > 50 {
				break
			}

			k := arg.Name
			if k == "" {
				k = cast.ToString(arg.Ordinal)
			}

			data[k] = cast.ToString(arg.Value)
		}
		s.Data = data
	})
	span.Finish()
}

// ESConn returns established connection
func ESConn(env Environment) (*elasticsearch.TypedClient, error) {
	return esConn(env, 30*time.Second)
}

// ESBulkConn returns established connection
func ESBulkConn(env Environment) (*elasticsearch.TypedClient, error) {
	return esConn(env, 360*time.Second)
}

func esConn(env Environment, timeout time.Duration) (*elasticsearch.TypedClient, error) {
	// Same retries as the previous client: 6 retries with fixed delays; a later attempt gets no delay.
	backoffs := []time.Duration{200 * time.Millisecond, 300 * time.Millisecond, 400 * time.Millisecond, 600 * time.Millisecond, 700 * time.Millisecond, 800 * time.Millisecond}

	// Start from the default transport (30s dial, idle-connection reaping) and bound the header wait only.
	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("uninitialized es client <%s>: http.DefaultTransport is %T", env.EnvString("ESURL"), http.DefaultTransport)
	}
	transport := base.Clone()
	transport.ResponseHeaderTimeout = timeout

	// No node discovery: every request goes through the Service in front of
	// the cluster, so sniffed pod addresses cannot outlive their pods.
	opts := []elasticsearch.Option{
		elasticsearch.WithAddresses(env.EnvString("ESURL")),
		elasticsearch.WithTransportOptions(
			elastictransport.WithTransport(transport),
			elastictransport.WithMaxRetries(len(backoffs)),
			elastictransport.WithRetryBackoff(func(attempt int) time.Duration {
				if attempt < 1 || attempt > len(backoffs) {
					return 0
				}
				return backoffs[attempt-1]
			}),
		),
	}
	if env.IsDebug() {
		opts = append(opts, elasticsearch.WithLogger(&elastictransport.TextLogger{Output: os.Stderr, EnableRequestBody: true, EnableResponseBody: true}))
	} else {
		opts = append(opts, elasticsearch.WithLogger(&logger.SentryErrorLogger{}))
	}

	es, err := elasticsearch.NewTyped(opts...)
	if err != nil {
		return nil, fmt.Errorf("uninitialized es client <%s>: %s", env.EnvString("ESURL"), err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	info, err := es.Info().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error got es version <%s>: %s", env.EnvString("ESURL"), err)
	}

	msg := "[INFO] the elasticsearch connection established <%s>, version %s"
	logger.Printf(msg, env.EnvString("ESURL"), info.Version.Int)
	return es, nil
}

// RedisConn returns established connection
func RedisConn(env Environment) (*radix.Pool, error) {
	return SelectRedisConn(env.EnvString("RedisURI"))
}

// SelectRedisConn returns established connection
func SelectRedisConn(uri string) (*radix.Pool, error) {
	dr, err := dsn.Redis(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis dsn <%s>: %s", uri, err)
	}

	selectDB, err := strconv.Atoi(dr.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis db number <%s>: %s", uri, err)
	}

	// this is a ConnFunc which will set up a connection which is authenticated
	// and has a 1 minute timeout on all operations
	connFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(time.Second*10),
			radix.DialSelectDB(selectDB),
		)
	}

	p, err := radix.NewPool("tcp", dr.HostPort, 10, radix.PoolConnFunc(connFunc))
	if err != nil {
		return nil, fmt.Errorf("uninitialized redis client <%s>: %s", uri, err)
	}

	msg := "[INFO] the redis@v3 connection established <%s>, version UNKNOWN"
	logger.Printf(msg, uri)

	return p, err
}

// DLMConn returns distributed lock manager pool
func DLMConn(env Environment) (*dlm.DLM, error) {
	dr, err := dsn.Redis(env.EnvString("DLMURI"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DLM dsn <%s>: %s", env.EnvString("DLMURI"), err)
	}

	pool := &dlmredis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (dlmredis.Conn, error) {
			c, err := dlmredis.Dial("tcp", dr.HostPort)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", dr.DB); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c dlmredis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := dlmredis.String(conn.Do("PING")); err != nil {
		return nil, fmt.Errorf("uninitialized DLM client <%s>: %s", env.EnvString("DLMURI"), err)
	}

	msg := "[INFO] the DLM(distributed lock) connection established <%s>, version UNKNOWN"
	logger.Printf(msg, env.EnvString("DLMURI"))

	return &dlm.DLM{Pool: pool}, nil
}

// BQConn returns err
func BQConn(env Environment) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, env.EnvString("GCProject"))
	if err != nil {
		return fmt.Errorf("there is no project in bigquery <%s>: %s", env.EnvString("GCProject"), err)
	}
	defer client.Close()

	msg := "[INFO] the bigquery connection established <%s>"
	logger.Printf(msg, env.EnvString("GCProject"))
	return nil
}
