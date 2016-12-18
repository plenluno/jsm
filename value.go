package jsm

import (
	"math"
	"reflect"
)

// ToBool converts the given value to a boolean value.
func ToBool(v interface{}) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		f := val.Float()
		return f != 0.0 && !math.IsNaN(f)
	case reflect.String:
		return val.String() != ""
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return !val.IsNil()
	default:
		return true
	}
}
