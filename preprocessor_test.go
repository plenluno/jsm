package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocess(t *testing.T) {
	assert := assert.New(t)

	p1 := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{3, 4.5}},
		{Mnemonic: MnemonicJump, Immediates: []interface{}{"abc"}, Comment: "jump"},
		{Label: "abc", Mnemonic: MnemonicReturn, Immediates: []interface{}{2.0}},
	}
	p2 := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{3, 4.5}},
		{Mnemonic: MnemonicJump, Immediates: []interface{}{2}},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{2}},
	}

	assert.Equal(p2, preprocess(p1))
}
