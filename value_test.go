package jsm

import (
	"encoding/json"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(TypeNull, TypeOf(NullValue()))
	assert.Equal(TypeNull, TypeOf(ArrayValue(nil)))
	assert.Equal(TypeNull, TypeOf(ObjectValue(nil)))
	assert.Equal(TypeNull, TypeOf(PointerValue(nil)))
	assert.Equal(TypeBoolean, TypeOf(BooleanValue(false)))
	assert.Equal(TypeBoolean, TypeOf(BooleanValue(true)))
	assert.Equal(TypeNumber, TypeOf(IntegerValue(1)))
	assert.Equal(TypeNumber, TypeOf(NumberValue(1.1)))
	assert.Equal(TypeNumber, TypeOf(NumberValue(math.NaN())))
	assert.Equal(TypeNumber, TypeOf(NumberValue(math.Inf(1))))
	assert.Equal(TypeNumber, TypeOf(NumberValue(math.Inf(-1))))
	assert.Equal(TypeString, TypeOf(StringValue("a")))
	assert.Equal(TypeArray, TypeOf(ArrayValue([]Value{IntegerValue(123), StringValue("abc")})))
	assert.Equal(TypeObject, TypeOf(ObjectValue(map[string]Value{"abc": NumberValue(1.23)})))
	assert.Equal(TypePointer, TypeOf(PointerValue(unsafe.Pointer(assert))))
}

func TestToBoolean(t *testing.T) {
	assert := assert.New(t)

	assert.False(ToBoolean(NullValue()))
	assert.False(ToBoolean(ArrayValue(nil)))
	assert.False(ToBoolean(ObjectValue(nil)))
	assert.False(ToBoolean(PointerValue(nil)))
	assert.False(ToBoolean(BooleanValue(false)))
	assert.False(ToBoolean(IntegerValue(0)))
	assert.False(ToBoolean(NumberValue(0.0)))
	assert.False(ToBoolean(NumberValue(math.NaN())))
	assert.False(ToBoolean(StringValue("")))

	assert.True(ToBoolean(BooleanValue(true)))
	assert.True(ToBoolean(IntegerValue(1)))
	assert.True(ToBoolean(NumberValue(1.1)))
	assert.True(ToBoolean(NumberValue(math.Inf(1))))
	assert.True(ToBoolean(NumberValue(math.Inf(-1))))
	assert.True(ToBoolean(StringValue("a")))
	assert.True(ToBoolean(ArrayValue([]Value{})))
	assert.True(ToBoolean(ObjectValue(map[string]Value{})))
	assert.True(ToBoolean(PointerValue(unsafe.Pointer(assert))))
}

func TestToNumber(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0.0, ToNumber(NullValue()))
	assert.Equal(0.0, ToNumber(ArrayValue(nil)))
	assert.Equal(0.0, ToNumber(ObjectValue(nil)))
	assert.Equal(0.0, ToNumber(PointerValue(nil)))
	assert.Equal(0.0, ToNumber(BooleanValue(false)))
	assert.Equal(1.0, ToNumber(BooleanValue(true)))
	assert.Equal(1.0, ToNumber(IntegerValue(1)))
	assert.Equal(1.1, ToNumber(NumberValue(1.1)))
	assert.Equal(1.1, ToNumber(StringValue("1.1")))

	assert.True(math.IsNaN(ToNumber(NumberValue(math.NaN()))))
	assert.True(math.IsNaN(ToNumber(StringValue(""))))
	assert.True(math.IsNaN(ToNumber(StringValue("a"))))
	assert.True(math.IsNaN(ToNumber(StringValue("NaN"))))
	assert.True(math.IsNaN(ToNumber(ArrayValue([]Value{}))))
	assert.True(math.IsNaN(ToNumber(ObjectValue(map[string]Value{}))))
	assert.True(math.IsNaN(ToNumber(PointerValue(unsafe.Pointer(assert)))))

	assert.True(math.IsInf(ToNumber(NumberValue(math.Inf(1))), 1))
	assert.True(math.IsInf(ToNumber(NumberValue(math.Inf(-1))), -1))
	assert.True(math.IsInf(ToNumber(StringValue("Infinity")), 1))
	assert.True(math.IsInf(ToNumber(StringValue("-Infinity")), -1))
}

func TestToInteger(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, ToInteger(NullValue()))
	assert.Equal(0, ToInteger(ArrayValue(nil)))
	assert.Equal(0, ToInteger(ObjectValue(nil)))
	assert.Equal(0, ToInteger(PointerValue(nil)))
	assert.Equal(0, ToInteger(BooleanValue(false)))
	assert.Equal(1, ToInteger(BooleanValue(true)))
	assert.Equal(1, ToInteger(IntegerValue(1)))
	assert.Equal(1, ToInteger(NumberValue(1.1)))
	assert.Equal(-1, ToInteger(NumberValue(-1.1)))
	assert.Equal(1, ToInteger(StringValue("1.1")))
	assert.Equal(-1, ToInteger(StringValue("-1.1")))

	assert.Equal(0, ToInteger(NumberValue(math.NaN())))
	assert.Equal(0, ToInteger(StringValue("")))
	assert.Equal(0, ToInteger(StringValue("a")))
	assert.Equal(0, ToInteger(StringValue("NaN")))
	assert.Equal(0, ToInteger(ArrayValue([]Value{})))
	assert.Equal(0, ToInteger(ObjectValue(map[string]Value{})))
	assert.Equal(0, ToInteger(PointerValue(unsafe.Pointer(assert))))

	assert.Equal(maxInt, ToInteger(NumberValue(math.Inf(1))))
	assert.Equal(minInt, ToInteger(NumberValue(math.Inf(-1))))
	assert.Equal(maxInt, ToInteger(StringValue("Infinity")))
	assert.Equal(minInt, ToInteger(StringValue("-Infinity")))
}

func TestToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("null", ToString(NullValue()))
	assert.Equal("null", ToString(ArrayValue(nil)))
	assert.Equal("null", ToString(ObjectValue(nil)))
	assert.Equal("null", ToString(PointerValue(nil)))
	assert.Equal("false", ToString(BooleanValue(false)))
	assert.Equal("true", ToString(BooleanValue(true)))
	assert.Equal("1", ToString(IntegerValue(1)))
	assert.Equal("1", ToString(NumberValue(1.0)))
	assert.Equal("-1.1", ToString(NumberValue(-1.1)))
	assert.Equal("NaN", ToString(NumberValue(math.NaN())))
	assert.Equal("Infinity", ToString(NumberValue(math.Inf(1))))
	assert.Equal("-Infinity", ToString(NumberValue(math.Inf(-1))))
	assert.Equal("", ToString(StringValue("")))
	assert.Equal("a", ToString(StringValue("a")))
	assert.Equal("", ToString(ArrayValue([]Value{})))
	assert.Equal("123,abc", ToString(ArrayValue([]Value{IntegerValue(123), StringValue("abc")})))
	assert.Equal("{}", ToString(ObjectValue(map[string]Value{})))
	assert.Equal("{\"abc\":1.23}", ToString(ObjectValue(map[string]Value{"abc": NumberValue(1.23)})))
	assert.Equal("null", ToString(PointerValue(unsafe.Pointer(assert))))
}

func TestToPointer(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(unsafe.Pointer(nil), ToPointer(NullValue()))
	assert.Equal(unsafe.Pointer(nil), ToPointer(ArrayValue(nil)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(ObjectValue(nil)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(PointerValue(nil)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(BooleanValue(false)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(BooleanValue(true)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(IntegerValue(1)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(NumberValue(1.0)))
	assert.Equal(unsafe.Pointer(nil), ToPointer(StringValue("a")))
	assert.Equal(unsafe.Pointer(nil), ToPointer(ArrayValue([]Value{})))
	assert.Equal(unsafe.Pointer(nil), ToPointer(ObjectValue(map[string]Value{})))
	assert.Equal(unsafe.Pointer(assert), ToPointer(PointerValue(unsafe.Pointer(assert))))
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	assert.True(Equal(NullValue(), NullValue()))
	assert.True(Equal(NullValue(), ArrayValue(nil)))
	assert.True(Equal(NullValue(), ObjectValue(nil)))
	assert.True(Equal(NullValue(), PointerValue(nil)))
	assert.True(Equal(BooleanValue(false), BooleanValue(false)))
	assert.True(Equal(BooleanValue(true), BooleanValue(true)))
	assert.True(Equal(IntegerValue(1), IntegerValue(1)))
	assert.True(Equal(NumberValue(1.1), NumberValue(1.1)))
	assert.True(Equal(NumberValue(math.Inf(1)), NumberValue(math.Inf(1))))
	assert.True(Equal(NumberValue(math.Inf(-1)), NumberValue(math.Inf(-1))))
	assert.True(Equal(IntegerValue(9007199254740991), NumberValue(9007199254740991.0)))
	assert.True(Equal(StringValue("a"), StringValue("a")))
	assert.True(Equal(
		ArrayValue([]Value{IntegerValue(123), StringValue("abc")}),
		ArrayValue([]Value{NumberValue(123.0), StringValue("abc")})))
	assert.True(Equal(
		ObjectValue(map[string]Value{"abc": IntegerValue(123)}),
		ObjectValue(map[string]Value{"abc": NumberValue(123.0)})))
	assert.True(Equal(
		PointerValue(unsafe.Pointer(assert)),
		PointerValue(unsafe.Pointer(assert))))

	assert.False(Equal(NullValue(), BooleanValue(false)))
	assert.False(Equal(BooleanValue(false), BooleanValue(true)))
	assert.False(Equal(IntegerValue(0), IntegerValue(1)))
	assert.False(Equal(NumberValue(1.0), NumberValue(1.1)))
	assert.False(Equal(NumberValue(math.NaN()), NumberValue(math.NaN())))
	assert.False(Equal(IntegerValue(9007199254740992), NumberValue(9007199254740992.0)))
	assert.False(Equal(StringValue("a"), StringValue("b")))
	assert.False(Equal(
		ArrayValue([]Value{IntegerValue(123), StringValue("abc")}),
		ArrayValue([]Value{NumberValue(1.23), StringValue("abc")})))
	assert.False(Equal(
		ObjectValue(map[string]Value{"abc": IntegerValue(123)}),
		ObjectValue(map[string]Value{"abc": NumberValue(1.23)})))
	assert.False(Equal(
		PointerValue(unsafe.Pointer(nil)),
		PointerValue(unsafe.Pointer(assert))))
}

func TestLess(t *testing.T) {
	assert := assert.New(t)

	assert.True(Less(BooleanValue(false), BooleanValue(true)))
	assert.True(Less(IntegerValue(9007199254740991), IntegerValue(9007199254740992)))
	assert.True(Less(IntegerValue(9007199254740992), IntegerValue(9007199254740993)))
	assert.True(Less(IntegerValue(-9007199254740992), IntegerValue(-9007199254740991)))
	assert.True(Less(IntegerValue(-9007199254740993), IntegerValue(-9007199254740992)))
	assert.True(Less(NumberValue(1.0), NumberValue(1.1)))
	assert.True(Less(NumberValue(math.Inf(-1)), NumberValue(math.Inf(1))))
	assert.True(Less(StringValue("a"), StringValue("b")))

	assert.False(Less(NullValue(), NullValue()))
	assert.False(Less(NullValue(), BooleanValue(false)))
	assert.False(Less(BooleanValue(false), NullValue()))
	assert.False(Less(BooleanValue(false), BooleanValue(false)))
	assert.False(Less(BooleanValue(true), BooleanValue(true)))
	assert.False(Less(BooleanValue(true), BooleanValue(false)))
	assert.False(Less(IntegerValue(9007199254740991), IntegerValue(9007199254740991)))
	assert.False(Less(IntegerValue(9007199254740992), IntegerValue(9007199254740991)))
	assert.False(Less(IntegerValue(9007199254740992), IntegerValue(9007199254740992)))
	assert.False(Less(IntegerValue(9007199254740993), IntegerValue(9007199254740992)))
	assert.False(Less(IntegerValue(-9007199254740991), IntegerValue(-9007199254740991)))
	assert.False(Less(IntegerValue(-9007199254740991), IntegerValue(-9007199254740992)))
	assert.False(Less(IntegerValue(-9007199254740992), IntegerValue(-9007199254740992)))
	assert.False(Less(IntegerValue(-9007199254740992), IntegerValue(-9007199254740993)))
	assert.False(Less(NumberValue(1.0), NumberValue(1.0)))
	assert.False(Less(NumberValue(1.1), NumberValue(1.0)))
	assert.False(Less(NumberValue(math.Inf(1)), NumberValue(math.Inf(1))))
	assert.False(Less(NumberValue(math.Inf(-1)), NumberValue(math.Inf(-1))))
	assert.False(Less(NumberValue(math.Inf(1)), NumberValue(math.Inf(-1))))
	assert.False(Less(StringValue("a"), StringValue("a")))
	assert.False(Less(StringValue("b"), StringValue("a")))
}

func TestMarshalPointer(t *testing.T) {
	assert := assert.New(t)

	data, err := json.Marshal(PointerValue(nil))
	assert.NoError(err)
	assert.Equal([]byte("null"), data)

	data, err = json.Marshal(PointerValue(unsafe.Pointer(assert)))
	assert.NoError(err)
	assert.Equal([]byte("null"), data)
}
