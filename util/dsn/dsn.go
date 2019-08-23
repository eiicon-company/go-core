package dsn

import "fmt"

func ef(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
