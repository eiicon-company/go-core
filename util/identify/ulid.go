package identify

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULIDNow returns ULID string
func ULIDNow() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0) // nolint:gosec
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
