package jsm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocess(t *testing.T) {
	assert := assert.New(t)

	pp := newPreprocessor()

	p1 := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{3, 4.5}},
		{Mnemonic: MnemonicPop, Immediates: []interface{}{"1.1"}},
		{Mnemonic: MnemonicCall, Immediates: []interface{}{"abc", 1.1}},
		{Label: "abc", Mnemonic: MnemonicReturn, Immediates: []interface{}{2.2}},
		{Mnemonic: MnemonicJump, Immediates: []interface{}{"abc"}, Comment: "jump"},
	}
	p2 := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{3, 4.5}},
		{Mnemonic: MnemonicPop, Immediates: []interface{}{1}},
		{Mnemonic: MnemonicCall, Immediates: []interface{}{3, 1}},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{2}},
		{Mnemonic: MnemonicJump, Immediates: []interface{}{3}},
	}

	j1, err := json.Marshal(p1)
	assert.NoError(err)

	p3, err := pp.preprocess(p1)
	assert.NoError(err)
	assert.Equal(p2, p3)

	j2, err := json.Marshal(p1)
	assert.NoError(err)
	assert.Equal(j1, j2)
}
