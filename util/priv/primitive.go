package priv

import (
	"errors"
	"reflect"
)

var (
	floatType               = reflect.TypeOf(float64(0))
	stringType              = reflect.TypeOf("")
	errUnexpectedNumberType = errors.New("Non-numeric type could not be converted")
)
