package slices

import (
	"reflect"

	"golang.org/x/xerrors"
)

// Interfaces returns slice interface from interface
func Interfaces(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, xerrors.New("InterfaceSlice() given a non-slice type")
	}

	r := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		r[i] = s.Index(i).Interface()
	}

	return r, nil
}
