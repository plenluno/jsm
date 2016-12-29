package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapLoadStore(t *testing.T) {
	assert := assert.New(t)

	h := NewHeap()
	v, err := h.Load("abc")
	assert.Error(err)
	assert.Nil(v)

	h.Store("abc", 123)
	v, err = h.Load("abc")
	assert.NoError(err)
	assert.Equal(123, v)

	h.Store("abc", "xyz")
	v, err = h.Load("abc")
	assert.NoError(err)
	assert.Equal("xyz", v)
}

func TestHeapClear(t *testing.T) {
	assert := assert.New(t)

	h := NewHeap()
	h.Store("xyz", []int{1, 2, 3})
	h.Clear()
	d, err := h.Dump()
	assert.NoError(err)
	assert.Equal("{}", string(d))
}

func TestHeapDumpRestore(t *testing.T) {
	assert := assert.New(t)

	h1 := NewHeap()
	h1.Store("xyz", []int{1, 2, 3})
	d1, err := h1.Dump()
	assert.NoError(err)
	assert.Equal("{\"xyz\":[1,2,3]}", string(d1))

	h2 := NewHeap()
	d2, err := h2.Dump()
	assert.NoError(err)
	assert.Equal("{}", string(d2))

	err = h2.Restore(d1)
	assert.NoError(err)
	d2, err = h2.Dump()
	assert.NoError(err)
	assert.Equal("{\"xyz\":[1,2,3]}", string(d2))

	err = h2.Restore([]byte{})
	assert.Error(err)
}
