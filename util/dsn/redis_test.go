package dsn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	t.Helper()

	f, err := Redis("file://127.0.0.1:6379/4")
	require.Error(t, err)

	f, err = Redis("redis://127.0.0.1:6379/4")
	require.NoError(t, err)

	require.Empty(t, f.Password)
	require.Equal(t, "127.0.0.1", f.Host)
	require.Equal(t, "6379", f.Port)
	require.Equal(t, "127.0.0.1:6379", f.HostPort)
	require.Equal(t, "4", f.DB)

	t.Logf("Redis: %#+v", f)
}
