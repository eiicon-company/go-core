package dsn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	t.Helper()

	f, err := Redis("file://127.0.0.1:6379/4")
	require.Error(t, err, "Unknown Scheme: should return an error")

	f, err = Redis("redis://127.0.0.1:6379/4")
	require.NoError(t, err, "Unknown Scheme: should not return an error")

	require.Empty(t, f.Password, "redis field error")
	require.Equal(t, "127.0.0.1", f.Host, "redis field error")
	require.Equal(t, "6379", f.Port, "redis field error")
	require.Equal(t, "127.0.0.1:6379", f.HostPort, "redis field error")
	require.Equal(t, "4", f.DB, "redis field error")

	t.Logf("Redis: %#+v", f)
}
