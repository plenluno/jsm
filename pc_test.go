package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgramCounter(t *testing.T) {
	assert := assert.New(t)

	pc := newProgramCounter()
	assert.Equal(0, pc.Index())

	pc.SetIndex(7)
	assert.Equal(7, pc.Index())

	pc.Increment()
	assert.Equal(8, pc.Index())

	pc.Clear()
	assert.Equal(0, pc.Index())
}
