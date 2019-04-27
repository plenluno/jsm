// +build amd64

package jsm

import (
	"encoding/json"
	"math"
	"strconv"
	"unsafe"
)

const (
	bitsNaN     = 0x7ff0000000000001
	bitsNull    = 0x7ff1000000000000
	bitsTrue    = 0x7ff2000000000000
	bitsFalse   = 0xfff2000000000000
	bitsString  = 0x7ff3000000000000
	bitsArray   = 0x7ff4000000000000
	bitsObject  = 0x7ff5000000000000
	bitsPointer = 0x7ff6000000000000

	maskSign    = 0x8000000000000000
	maskNaN     = 0x7ff0000000000000
	maskTag     = 0x000f000000000000
	maskInteger = 0x0007ffffffffffff

	tagNaN       = 0x0000000000000000
	tagNull      = 0x0001000000000000
	tagBoolean   = 0x0002000000000000
	tagString    = 0x0003000000000000
	tagArray     = 0x0004000000000000
	tagObject    = 0x0005000000000000
	tagPointer   = 0x0006000000000000
	tagUndefined = 0x0007000000000000
	// if tag >= 0x0008000000000000, the value is an integer.
	tagInteger = 0x0008000000000000
)

// Value represents a JSON value.
type Value *value

type value struct {
	float   float64
	pointer unsafe.Pointer
}

var nullValue = &value{float: math.Float64frombits(bitsNull)}

// NullValue returns the null value.
func NullValue() Value {
	return nullValue
}

var (
	trueValue  = &value{float: math.Float64frombits(bitsTrue)}
	falseValue = &value{float: math.Float64frombits(bitsFalse)}
)

// BooleanValue returns the boolean value representing the specified bool.
func BooleanValue(b bool) Value {
	if b {
		return trueValue
	}
	return falseValue
}

const maxImmediate = 1<<51 - 1

var floatLargeInteger = math.Float64frombits(maskNaN | tagInteger | maskInteger)

// IntegerValue returns the integer value representing the specified int.
func IntegerValue(i int) Value {
	if i >= -maxImmediate && i <= maxImmediate {
		var bits uint64 = maskNaN | tagInteger
		if i >= 0 {
			bits |= uint64(i)
		} else {
			bits |= maskSign | uint64(i)&maskInteger
		}
		return &value{float: math.Float64frombits(bits), pointer: nil}
	}
	return &value{float: floatLargeInteger, pointer: unsafe.Pointer(&i)}
}

var (
	floatNaN = math.Float64frombits(bitsNaN)
	nanValue = &value{float: floatNaN}
)

// NumberValue returns the number value representing the specified float64.
func NumberValue(f float64) Value {
	if math.IsNaN(f) {
		return nanValue
	}
	return &value{float: f}
}

var floatString = math.Float64frombits(bitsString)

// StringValue returns the string value representing the specified string.
func StringValue(s string) Value {
	return &value{float: floatString, pointer: unsafe.Pointer(&s)}
}

var floatArray = math.Float64frombits(bitsArray)

// ArrayValue returns the array value representing the specified slice.
func ArrayValue(a []Value) Value {
	if a == nil {
		return nullValue
	}
	return &value{float: floatArray, pointer: unsafe.Pointer(&a)}
}

var floatObject = math.Float64frombits(bitsObject)

// ObjectValue returns the object value representing the specified map.
func ObjectValue(o map[string]Value) Value {
	if o == nil {
		return nullValue
	}
	return &value{float: floatObject, pointer: unsafe.Pointer(&o)}
}

var floatPointer = math.Float64frombits(bitsPointer)

// PointerValue returns the pointer value representing the specified unsafe.Pointer.
func PointerValue(p unsafe.Pointer) Value {
	if p == nil {
		return nullValue
	}
	return &value{float: floatPointer, pointer: p}
}

// TypeOf returns the type of the given value.
func TypeOf(v Value) Type {
	f := v.float
	if !math.IsNaN(f) {
		return TypeNumber
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return TypeNumber
	case tagNull:
		return TypeNull
	case tagBoolean:
		return TypeBoolean
	case tagString:
		return TypeString
	case tagArray:
		return TypeArray
	case tagObject:
		return TypeObject
	case tagPointer:
		return TypePointer
	case tagUndefined:
		return TypeUndefined
	default: // integer
		return TypeNumber
	}
}

// ToBoolean converts the given value to a boolean.
func ToBoolean(v Value) bool {
	f := v.float
	if !math.IsNaN(f) {
		return f != 0.0
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return false
	case tagNull:
		return false
	case tagBoolean:
		return bits == bitsTrue
	case tagString:
		return *(*string)(v.pointer) != ""
	case tagArray, tagObject, tagPointer:
		return true
	case tagUndefined:
		return false
	default: // integer
		return bits&maskInteger != 0
	}
}

// ToInteger converts the given value to an integer.
func ToInteger(v Value) int {
	f := v.float
	if !math.IsNaN(f) {
		return floatToInt(f)
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return 0
	case tagNull:
		return 0
	case tagBoolean:
		if bits == bitsTrue {
			return 1
		}
		return 0
	case tagString:
		f, err := strconv.ParseFloat(*(*string)(v.pointer), 64)
		if err != nil {
			return 0
		}
		return floatToInt(f)
	case tagArray, tagObject, tagPointer:
		return 0
	case tagUndefined:
		return 0
	default: // integer
		return toInteger(v)
	}
}

// ToNumber converts the given value to a floating point number.
func ToNumber(v Value) float64 {
	f := v.float
	if !math.IsNaN(f) {
		return f
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return floatNaN
	case tagNull:
		return 0.0
	case tagBoolean:
		if bits == bitsTrue {
			return 1.0
		}
		return 0.0
	case tagString:
		f, err := strconv.ParseFloat(*(*string)(v.pointer), 64)
		if err != nil {
			return floatNaN
		}
		return f
	case tagArray, tagObject, tagPointer:
		return floatNaN
	case tagUndefined:
		return floatNaN
	default: // integer
		return float64(toInteger(v))
	}
}

func toInteger(v Value) int {
	p := v.pointer
	if p != nil {
		return *(*int)(p)
	}

	f := v.float
	bits := math.Float64bits(f)
	if math.Signbit(f) {
		return -int(maskInteger - bits&maskInteger + 1)
	}
	return int(bits & maskInteger)
}

// ToString converts the given value to a string.
func ToString(v Value) string {
	f := v.float
	if !math.IsNaN(f) {
		return floatToString(f)
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return "NaN"
	case tagNull:
		return "null"
	case tagString:
		return *(*string)(v.pointer)
	case tagArray:
		a := *(*[]Value)(v.pointer)
		var s string
		for i, v := range a {
			if i != 0 {
				s += ","
			}
			s += ToString(v)
		}
		return s
	case tagPointer:
		return "null"
	default:
		data, err := json.Marshal(normalize(v))
		if err != nil {
			panic(err)
		}
		return string(data)
	}
}

// ToPointer converts the given value to a pointer.
func ToPointer(v Value) unsafe.Pointer {
	f := v.float
	if !math.IsNaN(f) {
		return unsafe.Pointer(nil)
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagPointer:
		return v.pointer
	default:
		return unsafe.Pointer(nil)
	}
}

// Equal checks if the given two values are equivalent.
func Equal(v1, v2 Value) bool {
	nv1 := normalize(v1)
	nv2 := normalize(v2)
	f1 := nv1.float
	f2 := nv2.float
	notNaN1 := !math.IsNaN(f1)
	notNaN2 := !math.IsNaN(f2)
	if notNaN1 && notNaN2 {
		return f1 == f2
	} else if notNaN1 || notNaN2 {
		return false
	}

	bits1 := math.Float64bits(f1)
	bits2 := math.Float64bits(f2)
	tag1 := bits1 & maskTag
	tag2 := bits2 & maskTag
	if tag1 >= tagInteger {
		tag1 = tagInteger
	}
	if tag2 >= tagInteger {
		tag2 = tagInteger
	}
	if tag1 != tag2 {
		return false
	}

	switch tag1 {
	case tagNaN:
		return false
	case tagNull:
		return true
	case tagBoolean:
		return bits1 == bits2
	case tagString:
		return *(*string)(nv1.pointer) == *(*string)(nv2.pointer)
	case tagArray:
		a1 := *(*[]Value)(nv1.pointer)
		a2 := *(*[]Value)(nv2.pointer)
		for i, v1 := range a1 {
			v2 := a2[i]
			if !Equal(v1, v2) {
				return false
			}
		}
		return true
	case tagObject:
		m1 := *(*map[string]Value)(nv1.pointer)
		m2 := *(*map[string]Value)(nv2.pointer)
		for mk, mv1 := range m1 {
			mv2 := m2[mk]
			if !Equal(mv1, mv2) {
				return false
			}
		}
		return true
	case tagPointer:
		return nv1.pointer == nv2.pointer
	case tagUndefined:
		return true
	default: // integer
		return toInteger(nv1) == toInteger(nv2)
	}
}

// Less checks if v1 is less than v2.
func Less(v1, v2 Value) bool {
	nv1 := normalize(v1)
	nv2 := normalize(v2)
	f1 := nv1.float
	f2 := nv2.float
	notNaN1 := !math.IsNaN(f1)
	notNaN2 := !math.IsNaN(f2)
	if notNaN1 && notNaN2 {
		return f1 < f2
	} else if notNaN1 || notNaN2 {
		return false
	}

	return false
	/*
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
	*/
}

const maxSafeInteger = 1<<53 - 1

func normalize(v Value) Value {
	f := v.float
	if !math.IsNaN(f) {
		return v
	}

	bits := math.Float64bits(f)
	tag := bits & maskTag
	switch {
	case tag >= tagInteger:
		i := toInteger(v)
		if i >= -maxSafeInteger && i <= maxSafeInteger {
			return NumberValue(float64(i))
		}
		return v
	case tag == tagArray:
		a := *(*[]Value)(v.pointer)
		n := len(a)
		na := make([]Value, n)
		for i := 0; i < n; i++ {
			na[i] = normalize(a[i])
		}
		return ArrayValue(na)
	case tag == tagObject:
		m := *(*map[string]Value)(v.pointer)
		nm := map[string]Value{}
		for mk, mv := range m {
			nm[mk] = normalize(mv)
		}
		return ObjectValue(nm)
	default:
		return v
	}
}

func (v *value) MarshalJSON() ([]byte, error) {
	f := v.float
	if !math.IsNaN(f) {
		return json.Marshal(f)
	}

	bits := math.Float64bits(f)
	switch bits & maskTag {
	case tagNaN:
		return []byte("null"), nil
	case tagNull:
		return []byte("null"), nil
	case tagBoolean:
		if bits == bitsTrue {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case tagString:
		return []byte("\"" + *(*string)(v.pointer) + "\""), nil
	case tagArray:
		a := *(*[]Value)(v.pointer)
		s := "["
		for i, v := range a {
			if i != 0 {
				s += ","
			}
			data, err := ((*value)(v)).MarshalJSON()
			if err != nil {
				return nil, err
			}
			s += string(data)
		}
		s += "]"
		return []byte(s), nil
	case tagObject:
		o := *(*map[string]Value)(v.pointer)
		s := "{"
		for k, v := range o {
			if len(s) > 1 {
				s += ","
			}
			s += "\"" + k + "\":"
			data, err := ((*value)(v)).MarshalJSON()
			if err != nil {
				return nil, err
			}
			s += string(data)
		}
		s += "}"
		return []byte(s), nil
	case tagPointer:
		return []byte("null"), nil
	case tagUndefined:
		return nil, nil
	default: // integer
		return json.Marshal(toInteger(v))
	}
}
