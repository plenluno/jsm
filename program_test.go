package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgram(t *testing.T) {
	assert := assert.New(t)

	p := Program{
		[]Instruction{
			{Mnemonic: MnemonicPush, Operands: []interface{}{3.0}},
			{Mnemonic: MnemonicPop},
		},
	}
	s := "{\"program\":[{\"mnemonic\":\"PUSH\",\"operands\":[3]},{\"mnemonic\":\"POP\"}]}"

	j, err := json.Marshal(p)
	assert.NoError(err)
	assert.Equal(s, string(j))

	var v Program
	err = json.Unmarshal([]byte(s), &v)
	assert.NoError(err)
	assert.Equal(p, v)
}
