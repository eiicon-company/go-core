package identify

import (
	"crypto/rand"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

// ULIDNow returns ULID string
func ULIDNow() string {
	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), rand.Reader).String()
}
