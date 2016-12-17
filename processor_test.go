package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTrue(t *testing.T) {
	assert := assert.New(t)

	type st struct {
		i int
	}
	var s st

	var p *int
	var sp *st
	var m map[int]int
	var sl []int

	assert.False(isTrue(nil))
	assert.False(isTrue(false))
	assert.False(isTrue(0))
	assert.False(isTrue(0.0))
	assert.False(isTrue(""))
	assert.False(isTrue(p))
	assert.False(isTrue(sp))
	assert.False(isTrue(m))
	assert.False(isTrue(sl))

	p = &s.i
	sp = &s
	m = map[int]int{}
	sl = []int{}

	assert.True(isTrue(true))
	assert.True(isTrue(1))
	assert.True(isTrue(1.1))
	assert.True(isTrue("a"))
	assert.True(isTrue(s))
	assert.True(isTrue(p))
	assert.True(isTrue(sp))
	assert.True(isTrue(m))
	assert.True(isTrue(sl))
}

func TestIsFalse(t *testing.T) {
	assert := assert.New(t)

	type st struct {
		i int
	}
	var s st

	var p *int
	var sp *st
	var m map[int]int
	var sl []int

	assert.True(isFalse(nil))
	assert.True(isFalse(false))
	assert.True(isFalse(0))
	assert.True(isFalse(0.0))
	assert.True(isFalse(""))
	assert.True(isFalse(p))
	assert.True(isFalse(sp))
	assert.True(isFalse(m))
	assert.True(isFalse(sl))

	p = &s.i
	sp = &s
	m = map[int]int{}
	sl = []int{}

	assert.False(isFalse(true))
	assert.False(isFalse(1))
	assert.False(isFalse(1.1))
	assert.False(isFalse("a"))
	assert.False(isFalse(s))
	assert.False(isFalse(p))
	assert.False(isFalse(sp))
	assert.False(isFalse(m))
	assert.False(isFalse(sl))
}
