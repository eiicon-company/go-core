package dlm

import (
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

type (
	// DLM is called distributed lock manager.
	DLM struct {
		Pool *redis.Pool
	}
)

// Mutex returns MUTual EXclusion
func (d *DLM) Mutex(name string, expires time.Duration) *redsync.Mutex {
	rs := redsync.New([]redsync.Pool{d.Pool})
	return rs.NewMutex(name, redsync.SetExpiry(expires))
}

// MutexOptions returns MUTual EXclusion
func (d *DLM) MutexOptions(name string, options ...redsync.Option) *redsync.Mutex {
	rs := redsync.New([]redsync.Pool{d.Pool})
	return rs.NewMutex(name, options...)
}

// Close dlm connection pooling
func (d *DLM) Close() error {
	//
	// return d.Pool.Close() // never close cause of one pool is used right now
	//
	return nil
}

func (d *DLM) Exists(name string) (bool, error) {
	conn := d.Pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", name))
	if err != nil && err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return reply != "", nil
}
