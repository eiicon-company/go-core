package priv

import (
	"reflect"

	"golang.org/x/xerrors"
)

var (
	floatType               = reflect.TypeOf(float64(0))
	stringType              = reflect.TypeOf("")
	errUnexpectedNumberType = xerrors.New("Non-numeric type could not be converted")
)
