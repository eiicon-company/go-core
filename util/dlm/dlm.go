package dlm

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/hako/branca"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// DLM is distributed lock manager
type DLM struct {
	Client *redis.Client
}

// New returns a new DLM.
func (d *DLM) New() *redsync.Redsync {
	pool := goredis.NewPool(d.Client)
	return redsync.New(pool)
}

// Lock is lock
func (d *DLM) Lock(key string, expired time.Duration) (*redsync.Mutex, error) {
	rs := d.New()
	mutex := rs.NewMutex(key, redsync.WithExpiry(expired))
	return mutex, mutex.Lock()
}

// Unlock is unlock
func (d *DLM) Unlock(mutex *redsync.Mutex) (bool, error) {
	return mutex.Unlock()
}

// Branca is a token generator
func (d *DLM) Branca(key string, expired time.Duration) (string, error) {
	mutex, err := d.Lock(key, expired)
	if err != nil {
		return "", errors.Wrap(err, "faild to lock")
	}
	defer d.Unlock(mutex)

	br := branca.NewBranca("12345678901234567890123456789012")
	br.SetTTL(uint32(expired.Seconds()))
	return br.EncodeToString(key)
}
