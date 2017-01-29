package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapLoadStore(t *testing.T) {
	assert := assert.New(t)

	h := newHeap()
	v, err := h.Load("abc")
	assert.Error(err)
	assert.Nil(v)

	h.Store("abc", IntegerValue(123))
	v, err = h.Load("abc")
	assert.NoError(err)
	assert.Equal(123, ToInteger(v))

	h.Store("abc", StringValue("xyz"))
	v, err = h.Load("abc")
	assert.NoError(err)
	assert.Equal("xyz", ToString(v))
}

func TestHeapClear(t *testing.T) {
	assert := assert.New(t)

	h := newHeap()
	vs := ArrayValue([]Value{IntegerValue(1), NumberValue(2.0), StringValue("3")})
	h.Store("xyz", vs)
	h.Clear()
	d, err := h.Dump()
	assert.NoError(err)
	assert.Equal("{}", string(d))
}

func TestHeapDumpRestore(t *testing.T) {
	assert := assert.New(t)

	h1 := newHeap()
	vs := ArrayValue([]Value{IntegerValue(1), NumberValue(2.0), StringValue("3")})
	h1.Store("xyz", vs)
	d1, err := h1.Dump()
	assert.NoError(err)
	assert.Equal("{\"xyz\":[1,2,\"3\"]}", string(d1))

	h2 := newHeap()
	d2, err := h2.Dump()
	assert.NoError(err)
	assert.Equal("{}", string(d2))

	err = h2.Restore(d1)
	assert.NoError(err)
	d2, err = h2.Dump()
	assert.NoError(err)
	assert.Equal("{\"xyz\":[1,2,\"3\"]}", string(d2))

	err = h2.Restore([]byte{})
	assert.Error(err)
}
