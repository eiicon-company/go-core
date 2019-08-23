package util

import (
	"reflect"
	"time"

	"github.com/eiicon-company/go-utils/util/logger"
	"github.com/imdario/mergo"
)

type (
	timeTransfomer struct{}
)

func (t timeTransfomer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {

		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})

				if result[0].Bool() {
					dst.Set(src)
				}
			}

			return nil
		}
	}

	return nil
}

// Merge assembles []map[string]interface{} to dest
func Merge(dest interface{}, values ...interface{}) error {
	data := make(map[string]interface{})

	for _, value := range values {
		v, _ := reflections.Items(value)
		if err := mergo.Map(&data, v); err != nil {
			logger.E("merge.go Merge: %s", err)
		}
	}

	return mergo.Map(dest, data, mergo.WithTransformers(timeTransfomer{}))
}

// OverwriteMerge assembles []map[string]interface{} to dest
func OverwriteMerge(dest interface{}, values ...interface{}) error {
	data := make(map[string]interface{})

	for _, value := range values {
		v, _ := reflections.Items(value)
		if err := mergo.MapWithOverwrite(&data, v); err != nil {
			logger.E("merge.go OverwriteMerge: %s", err)
		}
	}

	return mergo.Map(dest, data, mergo.WithTransformers(timeTransfomer{}))
}
