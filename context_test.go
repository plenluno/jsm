package jsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineContext(t *testing.T) {
	assert := assert.New(t)

	mc := newMachineContext(newMachine())
	assert.Equal(mc.Value(keyPC), GetProgramCounter(mc))
	assert.Equal(mc.Value(keyHeap), GetGlobalHeap(mc))
	assert.Equal(mc.Value(keyStack), getCallStack(mc))

	setResult(mc, IntegerValue(3))
	assert.Equal(3, ToInteger(getResult(mc)))
	assert.Equal(3, ToInteger(*mc.Value(keyResult).(*Value)))
}

func TestProgramContext(t *testing.T) {
	assert := assert.New(t)

	pc := newProgramContext()
	assert.Equal(pc.Value(keyLabels), GetLabels(pc))

	setMnemonic(pc, MnemonicAdd)
	assert.Equal(string(MnemonicAdd), string(GetMnemonic(pc)))
	assert.Equal(string(MnemonicAdd), string(*pc.Value(keyMnemonic).(*Mnemonic)))
}
