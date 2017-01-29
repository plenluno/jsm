package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrame(t *testing.T) {
	assert := assert.New(t)

	f := newFrame()
	f.Arguments = []Value{NumberValue(3.5)}
	f.Locals.Store("abc", NumberValue(123.0))
	f.Operands.Push(NullValue())
	f.Operands.Push(StringValue("xyz"))
	f.ReturnTo = 16

	s := "{\"arguments\":[3.5],\"locals\":{\"abc\":123},\"operands\":[null,\"xyz\"],\"returnTo\":16}"

	j, err := json.Marshal(f)
	assert.NoError(err)
	assert.Equal(s, string(j))

	var v frame
	err = json.Unmarshal([]byte(s), &v)
	assert.NoError(err)
	assert.Equal(f, &v)
}
