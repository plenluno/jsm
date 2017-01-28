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
			Immediates: []Value{3, 4.5},
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []Value{"1.1"},
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []Value{"abc", 1.1},
		},
		{
			Label:      "abc",
			Mnemonic:   MnemonicReturn,
			Immediates: []Value{2.2},
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []Value{"abc"},
			Comment:    "jump",
		},
	}
	p2 := []Instruction{
		{
			Mnemonic:   MnemonicPush,
			Immediates: []Value{3, 4.5},
			opcode:     opcode(MnemonicPush),
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []Value{1},
			opcode:     opcode(MnemonicPop),
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []Value{3, 1},
			opcode:     opcode(MnemonicCall),
		},
		{
			Mnemonic:   MnemonicReturn,
			Immediates: []Value{2},
			opcode:     opcode(MnemonicReturn),
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []Value{3},
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
