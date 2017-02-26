package jsm

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Value represents a JSON value.
type Value interface{}

// NullValue returns the null value.
func NullValue() Value {
	return nil
}

var trueValue Value = true
var falseValue Value = false

// BooleanValue returns the boolean value representing the specified bool.
func BooleanValue(b bool) Value {
	// return an already allocated boolean value
	// and avoids new memory allocation
	if b {
		return trueValue
	}
	return falseValue
}

// IntegerValue returns the integer value representing the specified int.
func IntegerValue(i int) Value {
	return i
}

// NumberValue returns the number value representing the specified float64.
func NumberValue(f float64) Value {
	return f
}

// StringValue returns the string value representing the specified string.
func StringValue(s string) Value {
	return s
}

// ArrayValue returns the array value representing the specified slice.
func ArrayValue(a []Value) Value {
	return a
}

// ObjectValue returns the object value representing the specified map.
func ObjectValue(o map[string]Value) Value {
	return o
}

// Type represents the type of a JSON value.
type Type int

// These constants are the types of JSON values.
const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeArray
	TypeObject
)

// TypeOf returns the type of the given value.
func TypeOf(v Value) Type {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Invalid:
		return TypeNull
	case reflect.Bool:
		return TypeBoolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return TypeNumber
	case reflect.String:
		return TypeString
	case reflect.Slice:
		return TypeArray
	case reflect.Map:
		return TypeObject
	case reflect.Ptr:
		if val.IsNil() {
			return TypeNull
		}
		return TypeUndefined
	default:
		return TypeUndefined
	}
}

// ToBoolean converts the given value to a boolean.
func ToBoolean(v Value) bool {
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

// ToInteger converts the given value to an integer.
func ToInteger(v Value) int {
	f := ToNumber(v)
	if math.IsNaN(f) {
		return 0
	}

	var sign float64
	if math.Signbit(f) {
		sign = -1.0
	} else {
		sign = 1.0
	}
	return int(sign * math.Floor(math.Abs(f)))
}

// ToNumber converts the given value to a floating point number.
func ToNumber(v Value) float64 {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Invalid:
		return 0.0
	case reflect.Bool:
		if val.Bool() {
			return 1.0
		}
		return 0.0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.String:
		f, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return math.NaN()
		}
		return f
	case reflect.Slice, reflect.Map, reflect.Ptr:
		if val.IsNil() {
			return 0.0
		}
		return math.NaN()
	default:
		return math.NaN()
	}
}

// ToString converts the given value to a string.
func ToString(v Value) string {
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
func Equal(v1, v2 Value) bool {
	return reflect.DeepEqual(normalize(v1), normalize(v2))
}

// Less checks if v1 is less than v2.
func Less(v1, v2 Value) bool {
	val1 := reflect.ValueOf(normalize(v1))
	val2 := reflect.ValueOf(normalize(v2))
	kind1 := val1.Kind()
	kind2 := val2.Kind()
	switch kind1 {
	case reflect.Bool:
		if kind2 == reflect.Bool {
			return !val1.Bool() && val2.Bool()
		}
		return false
	case reflect.Int64:
		switch kind2 {
		case reflect.Int64:
			return val1.Int() < val2.Int()
		case reflect.Uint64:
			// val1 < 0 && val2 >= 0
			return true
		case reflect.Float64:
			return float64(val1.Int()) < val2.Float()
		default:
			return false
		}
	case reflect.Uint64:
		switch kind2 {
		case reflect.Int64:
			// val1 >= 0 && val2 < 0
			return false
		case reflect.Uint64:
			return val1.Uint() < val2.Uint()
		case reflect.Float64:
			return float64(val1.Uint()) < val2.Float()
		default:
			return false
		}
	case reflect.Float64:
		switch kind2 {
		case reflect.Int64:
			return val1.Float() < float64(val2.Int())
		case reflect.Uint64:
			return val1.Float() < float64(val2.Uint())
		case reflect.Float64:
			return val1.Float() < val2.Float()
		default:
			return false
		}
	case reflect.String:
		if kind2 == reflect.String {
			return strings.Compare(val1.String(), val2.String()) < 0
		}
		return false
	default:
		return false
	}
}

const maxSafeInteger = 1<<53 - 1

func normalize(v Value) Value {
	val := reflect.ValueOf(v)
	kind := val.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		if i >= -maxSafeInteger && i <= maxSafeInteger {
			return float64(i)
		} else if i >= 0 {
			return uint64(i)
		} else if kind != reflect.Int64 {
			return i
		}
		return v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := val.Uint()
		if u <= maxSafeInteger {
			return float64(u)
		} else if kind != reflect.Uint64 {
			return u
		}
		return v
	case reflect.Float32:
		return val.Float()
	case reflect.Float64:
		return v
	case reflect.Slice:
		if val.IsNil() {
			return nil
		}

		n := val.Len()
		a := make([]Value, n)
		for i := 0; i < n; i++ {
			a[i] = normalize(val.Index(i).Interface())
		}
		return a
	case reflect.Map:
		if val.IsNil() {
			return nil
		}

		m := map[string]Value{}
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
