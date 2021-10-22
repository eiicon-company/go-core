package priv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/eiicon-company/go-core/util/logger"
)

// MustFloat returns
func MustFloat(unk interface{}) float64 {
	v, err := ToFloat(unk)
	if err != nil {
		msg := "[WARN] colud not cast to float64"
		logger.Println(msg, err.Error())
	}
	return v
}

// ToFloat returns
func ToFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.Indirect(reflect.ValueOf(unk))
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		}
		if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		}
		msg := "can't convert %v to float64"
		return math.NaN(), fmt.Errorf(msg, v.Type())
	}
}
