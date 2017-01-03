package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocess(t *testing.T) {
	assert := assert.New(t)

	p1 := []Instruction{
		{
			Mnemonic:   MnemonicPush,
			Immediates: []interface{}{3, 4.5},
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []interface{}{"1.1"},
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []interface{}{"abc", 1.1},
		},
		{
			Label:      "abc",
			Mnemonic:   MnemonicReturn,
			Immediates: []interface{}{2.2},
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []interface{}{"abc"},
			Comment:    "jump",
		},
	}
	p2 := []Instruction{
		{
			Mnemonic:   MnemonicPush,
			Immediates: []interface{}{3, 4.5},
			opcode:     opcode(MnemonicPush),
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []interface{}{1},
			opcode:     opcode(MnemonicPop),
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []interface{}{3, 1},
			opcode:     opcode(MnemonicCall),
		},
		{
			Mnemonic:   MnemonicReturn,
			Immediates: []interface{}{2},
			opcode:     opcode(MnemonicReturn),
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []interface{}{3},
			opcode:     opcode(MnemonicJump),
		},
	}

	before, err := json.Marshal(p1)
	assert.NoError(err)

	p3, err := newPreprocessor().preprocess(p1)
	assert.NoError(err)
	assert.Equal(p2, p3)

	after, err := json.Marshal(p1)
	assert.NoError(err)
	assert.Equal(before, after)
}
