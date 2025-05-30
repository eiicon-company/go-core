package priv

import (
	"strconv"

	"github.com/eiicon-company/go-core/util/logger"
)

// MustInt64 returns numeric value as int64
func MustInt64(unk interface{}) int64 {
	v, err := ToInt64(unk)
	if err != nil {
		logger.Warnf("colud not cast to int64: %v", err)
	}
	return v
}

// ToInt64 returns numeric value as int64
func ToInt64(unk interface{}) (int64, error) {
	v, err := ToInt(unk)
	if err != nil {
		logger.Warnf("could not cast to int64: %v", err)
	}

	return int64(v), nil
}

// MustInt returns
func MustInt(unk interface{}) int {
	v, err := ToInt(unk)
	if err != nil {
		logger.Warnf("colud not cast to int: %v", err)
	}
	return v
}

// ToInt returns
func ToInt(unk interface{}) (int, error) {
	switch i := unk.(type) {
	case float64:
		return int(i), nil
	case float32:
		return int(i), nil
	case int64:
		return int(i), nil
	case int32:
		return int(i), nil
	case int:
		return i, nil
	case uint64:
		if i > uint64(^uint(0)) {
			return 0, errUnexpectedNumberType // Handle overflow
		}
		//nolint:gosec // G115
		return int(i), nil
	case uint32:
		return int(i), nil
	case uint:
		if i > ^uint(0) {
			return 0, errUnexpectedNumberType // Handle overflow
		}
		//nolint:gosec // G115
		return int(i), nil
	case string:
		return strconv.Atoi(i)
	default:
		return 0, errUnexpectedNumberType
	}
}
