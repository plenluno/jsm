package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstruction(t *testing.T) {
	assert := assert.New(t)

	p := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{3.0}},
		{Mnemonic: MnemonicPop},
	}
	s := "[{\"mnemonic\":\"push\",\"immediates\":[3]},{\"mnemonic\":\"pop\"}]"

	j, err := json.Marshal(p)
	assert.NoError(err)
	assert.Equal(s, string(j))

	var v []Instruction
	err = json.Unmarshal([]byte(s), &v)
	assert.NoError(err)
	assert.Equal(p, v)
}
