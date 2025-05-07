package dsn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMailStdOut(t *testing.T) {
	f, err := Mail("stdout://")
	require.NoError(t, err)
	require.True(t, f.StdOut)
}

func TestMail(t *testing.T) {
	t.Helper()

	f, err := Mail("smtp://username@gmail.com:password@smtp.gmail.com(smtp.gmail.com:587)/?tls=false")
	require.NoError(t, err)

	require.Equal(t, "username@gmail.com", f.User)
	require.Equal(t, "password", f.Password)
	require.Equal(t, "smtp.gmail.com", f.Host)
	require.Equal(t, "smtp.gmail.com:587", f.Addr)
	require.Equal(t, "smtp.gmail.com", f.TLSServer)
	require.False(t, f.TLS)
	require.False(t, f.StdOut)

	t.Logf("Mail: %#+v", f)
}
