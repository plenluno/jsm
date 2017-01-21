package jsm

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBoolean(t *testing.T) {
	assert := assert.New(t)

	type st struct {
		i int
	}
	var s st

	var p *int
	var sp *st
	var m map[int]int
	var sl []int

	assert.False(ToBoolean(nil))
	assert.False(ToBoolean(false))
	assert.False(ToBoolean(0))
	assert.False(ToBoolean(0.0))
	assert.False(ToBoolean(math.NaN()))
	assert.False(ToBoolean(""))
	assert.False(ToBoolean(p))
	assert.False(ToBoolean(sp))
	assert.False(ToBoolean(m))
	assert.False(ToBoolean(sl))

	p = &s.i
	sp = &s
	m = map[int]int{}
	sl = []int{}

	assert.True(ToBoolean(true))
	assert.True(ToBoolean(1))
	assert.True(ToBoolean(1.1))
	assert.True(ToBoolean("a"))
	assert.True(ToBoolean(s))
	assert.True(ToBoolean(p))
	assert.True(ToBoolean(sp))
	assert.True(ToBoolean(m))
	assert.True(ToBoolean(sl))
}

func TestToNumber(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0.0, ToNumber(nil))
	assert.Equal(0.0, ToNumber([]int(nil)))
	assert.Equal(0.0, ToNumber(false))
	assert.Equal(1.0, ToNumber(true))
	assert.Equal(1.0, ToNumber(1))
	assert.Equal(1.1, ToNumber(1.1))
	assert.Equal(1.1, ToNumber("1.1"))
	assert.True(math.IsNaN(ToNumber("")))
}

func TestToInteger(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, ToInteger(nil))
	assert.Equal(0, ToInteger([]int(nil)))
	assert.Equal(0, ToInteger(false))
	assert.Equal(1, ToInteger(true))
	assert.Equal(1, ToInteger(1))
	assert.Equal(1, ToInteger(1.1))
	assert.Equal(-1, ToInteger(-1.1))
	assert.Equal(1, ToInteger("1.1"))
	assert.Equal(0, ToInteger(""))
}

func TestToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("null", ToString(nil))
	assert.Equal("null", ToString([]int(nil)))
	assert.Equal("false", ToString(false))
	assert.Equal("-1", ToString(-1))
	assert.Equal("1", ToString(1.0))
	assert.Equal("a", ToString("a"))
	assert.Equal("1,2", ToString([]interface{}{1.0, 2.0}))
	assert.Equal("{\"1\":2}", ToString(map[float64]int{1.0: 2}))
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

func TestLess(t *testing.T) {
	assert := assert.New(t)

	assert.True(Less(false, true))
	assert.True(Less(1.0, 1.1))
	assert.True(Less(9007199254740991, 9007199254740992))
	assert.True(Less(9007199254740992, 9007199254740993))
	assert.True(Less(-9007199254740992, -9007199254740991))
	assert.True(Less(-9007199254740993, -9007199254740992))
	assert.True(Less(-9007199254740993, 9007199254740993))
	assert.True(Less("a", "b"))

	assert.False(Less(nil, false))
	assert.False(Less(false, nil))
	assert.False(Less(true, false))
	assert.False(Less(false, false))
	assert.False(Less(1.1, 1.0))
	assert.False(Less(1.1, 1.1))
	assert.False(Less(9007199254740992, 9007199254740991))
	assert.False(Less(9007199254740993, 9007199254740992))
	assert.False(Less(-9007199254740991, -9007199254740992))
	assert.False(Less(-9007199254740992, -9007199254740993))
	assert.False(Less(9007199254740993, -9007199254740993))
	assert.False(Less(9007199254740993, 9007199254740993.0))
	assert.False(Less("b", "a"))
	assert.False(Less("b", "b"))
}
