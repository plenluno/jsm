package jsm

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type machineContextKey int

const (
	keyPC machineContextKey = iota
	keyHeap
	keyStack
	keyResult
)

type machineContext map[machineContextKey]interface{}

func (mc *machineContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (mc *machineContext) Done() <-chan struct{} {
	return nil
}

func (mc *machineContext) Err() error {
	return nil
}

func (mc *machineContext) Value(key interface{}) interface{} {
	switch key.(type) {
	case machineContextKey:
		return (*mc)[key.(machineContextKey)]
	default:
		return nil
	}
}

func newMachineContext(m *machine) context.Context {
	return &machineContext{
		keyPC:     m.PC,
		keyHeap:   m.Heap,
		keyStack:  m.Stack,
		keyResult: new(Value),
	}
}

// GetProgramCounter retrieves the program counter.
func GetProgramCounter(ctx context.Context) ProgramCounter {
	return (*ctx.(*machineContext))[keyPC].(*programCounter)
}

// GetGlobalHeap retrieves the global heap.
func GetGlobalHeap(ctx context.Context) Heap {
	return (*ctx.(*machineContext))[keyHeap].(*heap)
}

func getCallStack(ctx context.Context) *callStack {
	return (*ctx.(*machineContext))[keyStack].(*callStack)
}

func getFrame(ctx context.Context) (*frame, error) {
	return getCallStack(ctx).Peek()
}

// GetArgument retrieves the argument at the specified position.
func GetArgument(ctx context.Context, idx int) (Value, error) {
	f, err := getFrame(ctx)
	if err != nil {
		return NullValue(), err
	}

	if idx >= len(f.Arguments) {
		return NullValue(), errors.New("argument out of range")
	}

	return f.Arguments[idx], nil
}

// GetLocalHeap retrieves the current local heap.
func GetLocalHeap(ctx context.Context) (Heap, error) {
	f, err := getFrame(ctx)
	if err != nil {
		return nil, err
	}
	return f.Locals, nil
}

// GetOperandStack retrieves the current operand stack.
func GetOperandStack(ctx context.Context) (Stack, error) {
	f, err := getFrame(ctx)
	if err != nil {
		return nil, err
	}
	return f.Operands, nil
}

func getResult(ctx context.Context) Value {
	return *(*ctx.(*machineContext))[keyResult].(*Value)
}

func setResult(ctx context.Context, res Value) {
	r := (*ctx.(*machineContext))[keyResult].(*Value)
	*r = res
}

type programContextKey int

const (
	keyLabels programContextKey = iota
	keyMnemonic
)

type programContext map[programContextKey]interface{}

func (pc *programContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (pc *programContext) Done() <-chan struct{} {
	return nil
}

func (pc *programContext) Err() error {
	return nil
}

func (pc *programContext) Value(key interface{}) interface{} {
	switch key.(type) {
	case programContextKey:
		return (*pc)[key.(programContextKey)]
	default:
		return nil
	}
}

func newProgramContext() context.Context {
	return &programContext{
		keyLabels:   map[string]int{},
		keyMnemonic: new(Mnemonic),
	}
}

// GetLabels retrieves the program labels.
func GetLabels(ctx context.Context) map[string]int {
	return (*ctx.(*programContext))[keyLabels].(map[string]int)
}

// GetMnemonic retrieves the currently preprocessed mnemonic.
func GetMnemonic(ctx context.Context) Mnemonic {
	return *(*ctx.(*programContext))[keyMnemonic].(*Mnemonic)
}

func setMnemonic(ctx context.Context, mnemonic Mnemonic) {
	m := (*ctx.(*programContext))[keyMnemonic].(*Mnemonic)
	*m = mnemonic
}
