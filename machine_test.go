package jsm

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineRunFib(t *testing.T) {
	assert := assert.New(t)

	j, err := ioutil.ReadFile("./examples/fibonacci.json")
	assert.NoError(err)

	var p []Instruction
	err = json.Unmarshal(j, &p)
	assert.NoError(err)

	m := NewMachine()
	res, err := m.Run(p, []interface{}{1.0})
	assert.NoError(err)
	assert.Equal([]interface{}{1.0}, res)

	res, err = m.Run(p, []interface{}{6.0})
	assert.NoError(err)
	assert.Equal([]interface{}{13.0}, res)
}

func TestMachineRunSum(t *testing.T) {
	assert := assert.New(t)

	j, err := ioutil.ReadFile("./examples/sum_of_series.json")
	assert.NoError(err)

	var p []Instruction
	err = json.Unmarshal(j, &p)
	assert.NoError(err)

	m := NewMachine()
	res, err := m.Run(p, []interface{}{0.0})
	assert.NoError(err)
	assert.Equal([]interface{}{0.0}, res)

	res, err = m.Run(p, []interface{}{10.0})
	assert.NoError(err)
	assert.Equal([]interface{}{55.0}, res)
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

	GetProgramCounter(ctx).Increment()
	return nil
}

func TestMachineExtendFib(t *testing.T) {
	assert := assert.New(t)

	m := NewMachine()
	err := m.Extend("fib", fib, nil)
	assert.NoError(err)

	p := []Instruction{
		{Mnemonic: MnemonicLoadArgument, Immediates: []interface{}{0}},
		{Mnemonic: "fib"},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{1}},
	}
	res, err := m.Run(p, []interface{}{6})
	assert.NoError(err)
	assert.Equal([]interface{}{13}, res)
}

func sumOfSeries(n int) int {
	var sum int
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

func sum(ctx context.Context, imms []interface{}) error {
	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	if err := Push(ctx, sumOfSeries(ToInteger(v))); err != nil {
		return err
	}

	GetProgramCounter(ctx).Increment()
	return nil
}

func TestMachineExtendSum(t *testing.T) {
	assert := assert.New(t)

	m := NewMachine()
	err := m.Extend("sum", sum, nil)
	assert.NoError(err)

	p := []Instruction{
		{Mnemonic: MnemonicLoadArgument, Immediates: []interface{}{0}},
		{Mnemonic: "sum"},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{1}},
	}
	res, err := m.Run(p, []interface{}{10})
	assert.NoError(err)
	assert.Equal([]interface{}{55}, res)
}

func BenchmarkFibJSM(b *testing.B) {
	m := NewMachine()

	var p []Instruction
	j, _ := ioutil.ReadFile("./examples/fibonacci.json")
	json.Unmarshal(j, &p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Run(p, []interface{}{20.0})
	}
}

func BenchmarkFibNative(b *testing.B) {
	m := NewMachine()
	m.Extend("fib", fib, nil)

	p := []Instruction{
		{Mnemonic: MnemonicLoadArgument, Immediates: []interface{}{0}},
		{Mnemonic: "fib"},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{1}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Run(p, []interface{}{20})
	}
}

func BenchmarkSumJSM(b *testing.B) {
	m := NewMachine()

	var p []Instruction
	j, _ := ioutil.ReadFile("./examples/sum_of_series.json")
	json.Unmarshal(j, &p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Run(p, []interface{}{100000.0})
	}
}

func BenchmarkSumNative(b *testing.B) {
	m := NewMachine()
	m.Extend("sum", sum, nil)

	p := []Instruction{
		{Mnemonic: MnemonicLoadArgument, Immediates: []interface{}{0}},
		{Mnemonic: "sum"},
		{Mnemonic: MnemonicReturn, Immediates: []interface{}{1}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Run(p, []interface{}{100000})
	}
}
