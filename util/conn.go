package util

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/getsentry/sentry-go"
	"github.com/go-sql-driver/mysql"
	dlmredis "github.com/gomodule/redigo/redis"
	radix "github.com/mediocregopher/radix/v3"
	"github.com/olivere/elastic/v7"
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
	// https://github.com/go-sql-driver/mysql/issues/1302#issuecomment-1019842712
	db.SetConnMaxLifetime(time.Minute * 5) // https://github.com/go-sql-driver/mysql/issues/1120#issuecomment-636795680
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(8)

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
		PreExec: func(ctx context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostExec: func(ctx context.Context, dt interface{}, stmt *proxy.Stmt, args []driver.NamedValue, result driver.Result, err error) error {
			startTime := dt.(time.Time)
			since := time.Since(startTime)

			if since > period {
				span := sentry.StartSpan(ctx, "db.sql.exec.slow", func(s *sentry.Span) {
					s.StartTime = startTime
					s.EndTime = time.Now().Add(since)
					s.Description = stmt.QueryString

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

				ctx = span.Context() //nolint
			}

			return nil
		},
		PreQuery: func(ctx context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
			return time.Now(), nil
		},
		PostQuery: func(ctx context.Context, dt interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows, err error) error {
			startTime := dt.(time.Time)
			since := time.Since(startTime)

			if since > period {
				span := sentry.StartSpan(ctx, "db.sql.query.slow", func(s *sentry.Span) {
					s.StartTime = startTime
					s.EndTime = time.Now().Add(since)
					s.Description = stmt.QueryString

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

				ctx = span.Context() //nolint
			}

			return nil
		},
	}))
}

// ESConn returns established connection
func ESConn(env Environment) (*elastic.Client, error) {
	var op []elastic.ClientOptionFunc
	op = append(op, elastic.SetHttpClient(&http.Client{Timeout: 30 * time.Second}))
	op = append(op, elastic.SetURL(env.EnvString("ESURL")))
	op = append(op, elastic.SetSniff(true))
	op = append(op, elastic.SetHealthcheck(true))
	op = append(op, elastic.SetErrorLog(&logger.SentryErrorLogger{}))
	// 8 retries with fixed delay of 100ms, 200ms, 300ms, 400ms, 500ms, 600ms, 700ms, and 800ms.
	op = append(op, elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(100, 200, 300, 400, 600, 700, 800))))

	if env.IsDebug() {
		op = append(op, elastic.SetTraceLog(log.New(os.Stderr, "[[ELASTIC]] ", log.LstdFlags)))
		op = append(op, elastic.SetInfoLog(log.New(os.Stdout, "[ELASTIC] ", log.LstdFlags)))
	}

	es, err := elastic.NewClient(op...)
	if err != nil {
		return nil, fmt.Errorf("uninitialized es client <%s>: %s", env.EnvString("ESURL"), err)
	}
	ver, err := es.ElasticsearchVersion(env.EnvString("ESURL"))
	if err != nil {
		return nil, fmt.Errorf("error got es version <%s>: %s", env.EnvString("ESURL"), err)
	}

	msg := "[INFO] the elasticsearch connection established <%s>, version %s"
	logger.Printf(msg, env.EnvString("ESURL"), ver)
	return es, nil
}

// RedisConn returns established connection
func RedisConn(env Environment) (*radix.Pool, error) {
	uri := env.EnvString("RedisURI")

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
