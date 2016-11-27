package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert := assert.New(t)

	a := NewAddress()
	assert.Equal(0, a.Value())

	a.Increment()
	assert.Equal(1, a.Value())

	a.Jump(6)
	assert.Equal(7, a.Value())

	a.Jump(-4)
	assert.Equal(3, a.Value())

	a.Clear()
	assert.Equal(0, a.Value())
}
