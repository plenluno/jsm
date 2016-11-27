package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrame(t *testing.T) {
	assert := assert.New(t)

	f := NewFrame()
	f.Locals.Store("abc", 123.0)
	f.Operands.Push(123.0)
	f.Return.Address.Jump(16)

	s := "{\"locals\":{\"abc\":123},\"operands\":[123],\"return\":{\"address\":16}}"

	j, err := json.Marshal(f)
	assert.NoError(err)
	assert.Equal(s, string(j))

	var v Frame
	err = json.Unmarshal([]byte(s), &v)
	assert.NoError(err)
	assert.Equal(f, &v)
}
