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

	s.Push("abc")
	s.Push(123)

	v, err := s.Pop()
	assert.NoError(err)
	assert.Equal(123, v)

	v, err = s.Pop()
	assert.NoError(err)
	assert.Equal("abc", v)

	_, err = s.Pop()
	assert.Error(err)
}

func TestStackDo(t *testing.T) {
	assert := assert.New(t)

	s := newStack()
	s.Push(12)
	s.Push(3)
	s.Do(div, 2)
	v, err := s.Pop()
	assert.NoError(err)
	assert.Equal(4.0, v)

	_, err = s.Pop()
	assert.Error(err)
}

func TestStackPeekClear(t *testing.T) {
	assert := assert.New(t)

	s := newStack()
	s.Push("abc")
	v, err := s.Peek()
	assert.NoError(err)
	assert.Equal("abc", v)
	v, err = s.Peek()
	assert.NoError(err)
	assert.Equal("abc", v)

	s.Clear()
	_, err = s.Peek()
	assert.Error(err)
}

func TestStackDumpRestore(t *testing.T) {
	assert := assert.New(t)

	s1 := newStack()
	s1.Push("abc")
	s1.Push(123)
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
