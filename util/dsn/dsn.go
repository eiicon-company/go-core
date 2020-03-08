package dsn

import "fmt"

func ef(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// filePublicURL Http URL
var filePublicURL = "http://localhost:8000"
