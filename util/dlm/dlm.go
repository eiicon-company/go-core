package dlm

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/redis/go-redis/v9"
)

type (
	// DLM is called distributed lock manager.
	DLM struct {
		redsyncredis.Pool
	}
)

// Mutex returns MUTual EXclusion
func (d *DLM) Mutex(name string, expires time.Duration) *redsync.Mutex {
	rs := redsync.New(d.Pool)
	return rs.NewMutex(name, redsync.WithExpiry(expires))
}

// MutexOptions returns MUTual EXclusion
func (d *DLM) MutexOptions(name string, options ...redsync.Option) *redsync.Mutex {
	rs := redsync.New(d.Pool)
	return rs.NewMutex(name, options...)
}

// Close dlm connection pooling
// XXX: deprecated. do not need to close cause of one pool is used right now
func (d *DLM) Close() error {
	return nil
}

// Exists checks whether data exist or not.
func (d *DLM) Exists(ctx context.Context, name string) (bool, error) {
	conn, err := d.Pool.Get(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	reply, err := conn.Get(name)
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return reply != "", nil
}
