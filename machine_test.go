package jsm

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineRun(t *testing.T) {
	assert := assert.New(t)

	j, err := ioutil.ReadFile("./examples/fibonacci.json")
	assert.NoError(err)

	var p []Instruction
	err = json.Unmarshal(j, &p)
	assert.NoError(err)

	m := NewMachine()
	res, err := m.Run(p, []interface{}{1})
	assert.NoError(err)
	assert.Equal([]interface{}{1.0}, res)

	res, err = m.Run(p, []interface{}{6})
	assert.NoError(err)
	assert.Equal([]interface{}{13.0}, res)
}

func fibonacci(n int) int {
	if n < 2 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func fib(ctx context.Context, imms []interface{}) error {
	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	if err := Push(ctx, fibonacci(ToInteger(v))); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func TestMachineExtend(t *testing.T) {
	assert := assert.New(t)

	m := NewMachine()
	err := m.Extend("fib", fib, nil)
	assert.NoError(err)

	p := []Instruction{
		{Mnemonic: MnemonicPush, Immediates: []interface{}{0}},
		{Mnemonic: MnemonicLoadArgument},
		{Mnemonic: "fib"},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{1}},
	}
	res, err := m.Run(p, []interface{}{6})
	assert.NoError(err)
	assert.Equal([]interface{}{13}, res)
}
