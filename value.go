package jsm

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
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
	case reflect.Slice, reflect.Map, reflect.Ptr:
		return !val.IsNil()
	default:
		return true
	}
}

// ToString converts the given value to a string value.
func ToString(v interface{}) string {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Slice:
		if val.IsNil() {
			return "null"
		}

		n := val.Len()
		ss := make([]string, n)
		for i := 0; i < n; i++ {
			ss[i] = ToString(val.Index(i).Interface())
		}
		return strings.Join(ss, ",")
	default:
		data, err := json.Marshal(normalize(v))
		if err != nil {
			panic(err)
		}
		return string(data)
	}
}

// Equal checks if the given two values are equivalent.
func Equal(v1, v2 interface{}) bool {
	return reflect.DeepEqual(normalize(v1), normalize(v2))
}

const maxSafeInteger = 1<<53 - 1

func normalize(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if i >= -maxSafeInteger && i <= maxSafeInteger {
			return float64(i)
		} else if i >= 0 {
			return uint64(i)
		}
		return i
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := val.Uint()
		if u <= maxSafeInteger {
			return float64(u)
		}
		return u
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Slice:
		if val.IsNil() {
			return nil
		}

		n := val.Len()
		a := make([]interface{}, n)
		for i := 0; i < n; i++ {
			a[i] = normalize(val.Index(i).Interface())
		}
		return a
	case reflect.Map:
		if val.IsNil() {
			return nil
		}

		m := map[string]interface{}{}
		for _, key := range val.MapKeys() {
			m[ToString(key.Interface())] = normalize(val.MapIndex(key).Interface())
		}
		return m
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}

		return v
	default:
		return v
	}
}
