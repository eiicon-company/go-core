// Package scast is safety type casting
//
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package scast

import (
	"runtime"

	"github.com/lunemec/as"

	"github.com/eiicon-company/go-core/util/logger"
)

type numericType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Int returns ...
func Int[FromType numericType](value FromType) int {
	i, err := as.Int(value)
	if err != nil {
		_, f, l, ok := runtime.Caller(1)
		if !ok {
			return i
		}
		logger.C("%s:%d: Int overflow or underflow occurred %d: %+v", f, l, i, err)
	}
	return i
}

// Int8 returns ...
func Int8[FromType numericType](value FromType) int8 {
	i, err := as.Int8(value)
	if err != nil {
		_, f, l, ok := runtime.Caller(1)
		if !ok {
			return i
		}
		logger.C("%s:%d: Int8 overflow or underflow occurred %d: %+v", f, l, i, err)
	}
	return i
}

// Int32 returns ...
func Int32[FromType numericType](value FromType) int32 {
	i, err := as.Int32(value)
	if err != nil {
		_, f, l, ok := runtime.Caller(1)
		if !ok {
			return i
		}
		logger.C("%s:%d: Int32 overflow or underflow occurred %d: %+v", f, l, i, err)
	}
	return i
}

// Int64 returns ...
func Int64[FromType numericType](value FromType) int64 {
	i, err := as.Int64(value)
	if err != nil {
		_, f, l, ok := runtime.Caller(1)
		if !ok {
			return i
		}
		logger.C("%s:%d: Int64 overflow or underflow occurred %d: %+v", f, l, i, err)
	}
	return i
}
