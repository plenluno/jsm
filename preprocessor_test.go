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
			Immediates: []Value{IntegerValue(3), NumberValue(4.5)},
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []Value{StringValue("1.1")},
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []Value{StringValue("abc"), NumberValue(1.1)},
		},
		{
			Label:      "abc",
			Mnemonic:   MnemonicReturn,
			Immediates: []Value{NumberValue(2.2)},
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []Value{StringValue("abc")},
			Comment:    "jump",
		},
	}
	p2 := []Instruction{
		{
			Mnemonic:   MnemonicPush,
			Immediates: []Value{IntegerValue(3), NumberValue(4.5)},
			opcode:     opcode(MnemonicPush),
		},
		{
			Mnemonic:   MnemonicPop,
			Immediates: []Value{IntegerValue(1)},
			opcode:     opcode(MnemonicPop),
		},
		{
			Mnemonic:   MnemonicCall,
			Immediates: []Value{IntegerValue(3), IntegerValue(1)},
			opcode:     opcode(MnemonicCall),
		},
		{
			Mnemonic:   MnemonicReturn,
			Immediates: []Value{IntegerValue(2)},
			opcode:     opcode(MnemonicReturn),
		},
		{
			Mnemonic:   MnemonicJump,
			Immediates: []Value{IntegerValue(3)},
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
