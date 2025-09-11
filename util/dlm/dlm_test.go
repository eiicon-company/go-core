package dlm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRedisURI = "redis://localhost:6379/1" // Use a different DB for testing

// selectRedisConn is a simplified, local version of util.SelectRedisConn
// to avoid a circular dependency in tests.
func selectRedisConn(uri string) (*redis.Client, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis dsn <%s>: %w", uri, err)
	}

	opt.DialTimeout = time.Second * 10
	opt.MaxIdleConns = 10
	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("uninitialized redis client <%s>: %s", uri, err)
	}

	return rdb, nil
}

func setupDLM(t *testing.T) *DLM {
	t.Helper()

	rdb, err := selectRedisConn(testRedisURI)
	if err != nil {
		t.Fatalf("Failed to connect to Redis for testing. Is Redis running on %s? Error: %v", testRedisURI, err)
	}

	// Clean up keys after each test
	t.Cleanup(func() {
		rdb.Del(context.Background(), "test-lock")
		rdb.Del(context.Background(), "test-exists-lock")
		rdb.Close()
	})

	pool := goredis.NewPool(rdb)
	return &DLM{Pool: pool}
}

func TestDLM_Mutex(t *testing.T) {
	dlm := setupDLM(t)
	ctx := context.Background()

	mutex := dlm.Mutex("test-lock", 10*time.Second)

	// 1. Lock
	err := mutex.LockContext(ctx)
	require.NoError(t, err, "Failed to acquire lock")

	// 2. Check if lock exists
	exists, err := dlm.Exists(ctx, "test-lock")
	require.NoError(t, err)
	assert.True(t, exists, "Lock should exist after acquiring")

	// 3. Unlock
	ok, err := mutex.UnlockContext(ctx)
	require.NoError(t, err, "Failed to unlock")
	assert.True(t, ok, "Unlock should be successful")

	// 4. Check if lock still exists
	exists, err = dlm.Exists(ctx, "test-lock")
	require.NoError(t, err)
	assert.False(t, exists, "Lock should not exist after unlocking")

	// 5. Lock again to ensure it's reusable
	err = mutex.LockContext(ctx)
	require.NoError(t, err, "Failed to acquire lock again")

	// 6. Unlock again for cleanup
	_, err = mutex.UnlockContext(ctx)
	require.NoError(t, err)
}

func TestDLM_Exists(t *testing.T) {
	dlm := setupDLM(t)
	ctx := context.Background()
	lockKey := "test-exists-lock"

	// 1. Initially, the key should not exist
	exists, err := dlm.Exists(ctx, lockKey)
	require.NoError(t, err)
	assert.False(t, exists, "Lock should not exist initially")

	// 2. Acquire a lock
	mutex := dlm.Mutex(lockKey, 10*time.Second)
	err = mutex.LockContext(ctx)
	require.NoError(t, err)

	// 3. The key should now exist
	exists, err = dlm.Exists(ctx, lockKey)
	require.NoError(t, err)
	assert.True(t, exists, "Lock should exist after being acquired")

	// 4. Release the lock
	_, err = mutex.UnlockContext(ctx)
	require.NoError(t, err)

	// 5. The key should no longer exist
	exists, err = dlm.Exists(ctx, lockKey)
	require.NoError(t, err)
	assert.False(t, exists, "Lock should not exist after being released")
}
