package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackPushPop(t *testing.T) {
	assert := assert.New(t)

	s := newStack()
	_, err := s.Pop()
	assert.Error(err)

	s.Push(StringValue("abc"))
	s.Push(IntegerValue(123))

	v, err := s.Pop()
	assert.NoError(err)
	assert.Equal(123, ToInteger(v))

	v, err = s.Pop()
	assert.NoError(err)
	assert.Equal("abc", ToString(v))

	_, err = s.Pop()
	assert.Error(err)
}

func TestStackDo(t *testing.T) {
	assert := assert.New(t)

	s := newStack()
	s.Push(IntegerValue(12))
	s.Push(IntegerValue(3))
	s.Do(div, 2)
	v, err := s.Pop()
	assert.NoError(err)
	assert.Equal(4.0, ToNumber(v))

	_, err = s.Pop()
	assert.Error(err)
}

func TestStackPeekClear(t *testing.T) {
	assert := assert.New(t)

	s := newStack()
	s.Push(StringValue("abc"))
	v, err := s.Peek()
	assert.NoError(err)
	assert.Equal("abc", ToString(v))
	v, err = s.Peek()
	assert.NoError(err)
	assert.Equal("abc", ToString(v))

	s.Clear()
	_, err = s.Peek()
	assert.Error(err)
}

func TestStackDumpRestore(t *testing.T) {
	assert := assert.New(t)

	s1 := newStack()
	s1.Push(StringValue("abc"))
	s1.Push(IntegerValue(123))
	d1, err := s1.Dump()
	assert.NoError(err)
	assert.Equal("[\"abc\",123]", string(d1))

	s2 := newStack()
	d2, err := s2.Dump()
	assert.NoError(err)
	assert.Equal("[]", string(d2))

	err = s2.Restore(d1)
	assert.NoError(err)
	d2, err = s2.Dump()
	assert.NoError(err)
	assert.Equal("[\"abc\",123]", string(d2))

	err = s2.Restore([]byte{})
	assert.Error(err)
}
