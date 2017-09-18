package jsm

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToInt(t *testing.T) {
	assert := assert.New(t)

	assert.True(intSize == 32 || intSize == 64)

	if intSize == 32 {
		assert.Equal(math.MaxInt32, maxInt)
		assert.Equal(math.MinInt32, minInt)
	} else {
		assert.Equal(math.MaxInt64, maxInt)
		assert.Equal(math.MinInt64, minInt)
	}

	assert.Equal(0, floatToInt(math.NaN()))
	assert.Equal(maxInt, floatToInt(math.Inf(1)))
	assert.Equal(minInt, floatToInt(math.Inf(-1)))
}

func TestFloatToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("1", floatToString(1.0))
	assert.Equal("-1.1", floatToString(-1.1))
	assert.Equal("NaN", floatToString(math.NaN()))
	assert.Equal("Infinity", floatToString(math.Inf(1)))
	assert.Equal("-Infinity", floatToString(math.Inf(-1)))
}
