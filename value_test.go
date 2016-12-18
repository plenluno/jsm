package jsm

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert := assert.New(t)

	type st struct {
		i int
	}
	var s st

	var p *int
	var sp *st
	var m map[int]int
	var sl []int

	assert.False(ToBool(nil))
	assert.False(ToBool(false))
	assert.False(ToBool(0))
	assert.False(ToBool(0.0))
	assert.False(ToBool(math.NaN()))
	assert.False(ToBool(""))
	assert.False(ToBool(p))
	assert.False(ToBool(sp))
	assert.False(ToBool(m))
	assert.False(ToBool(sl))

	p = &s.i
	sp = &s
	m = map[int]int{}
	sl = []int{}

	assert.True(ToBool(true))
	assert.True(ToBool(1))
	assert.True(ToBool(1.1))
	assert.True(ToBool("a"))
	assert.True(ToBool(s))
	assert.True(ToBool(p))
	assert.True(ToBool(sp))
	assert.True(ToBool(m))
	assert.True(ToBool(sl))
}
