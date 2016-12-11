package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert := assert.New(t)

	a := NewAddress()
	assert.Equal(0, a.GetValue())

	a.SetValue(7)
	assert.Equal(7, a.GetValue())

	a.Increment()
	assert.Equal(8, a.GetValue())

	a.Clear()
	assert.Equal(0, a.GetValue())
}
