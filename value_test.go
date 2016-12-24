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

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	assert.True(Equal(nil, []int(nil)))
	assert.True(Equal(false, false))
	assert.True(Equal(0, 0.0))
	assert.True(Equal(int(1), uint(1)))
	assert.True(Equal(9007199254740991, 9007199254740991.0))
	assert.True(Equal("a", "a"))
	assert.True(Equal([]int{1, 2}, []interface{}{1.0, 2.0}))
	assert.True(Equal(map[float64]int{1.0: 2}, map[string]float64{"1": 2.0}))

	assert.False(Equal(nil, false))
	assert.False(Equal(true, false))
	assert.False(Equal(0, 1))
	assert.False(Equal(9007199254740992, 9007199254740992.0))
	assert.False(Equal("a", "b"))
	assert.False(Equal([]int{1, 2}, []interface{}{1, "a"}))
	assert.False(Equal(map[float64]int{1.0: 2}, map[string]float64{"1": 2.2}))
}
