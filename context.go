package jsm

import (
	"context"

	"github.com/pkg/errors"
)

type machineContextKey int

const (
	keyPC machineContextKey = iota
	keyHeap
	keyStack
	keyResult
)

func newMachineContext(m *machine) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyPC, Address(m.PC))
	ctx = context.WithValue(ctx, keyHeap, Heap(m.Heap))
	ctx = context.WithValue(ctx, keyStack, Stack(m.Stack))
	ctx = context.WithValue(ctx, keyResult, new(interface{}))
	return ctx
}

// GetPC retrieves the program counter from Context.
func GetPC(ctx context.Context) Address {
	return ctx.Value(keyPC).(Address)
}

// GetHeap retrieves Heap from Context.
func GetHeap(ctx context.Context) Heap {
	return ctx.Value(keyHeap).(Heap)
}

// GetStack retrieves Stack from Context.
func GetStack(ctx context.Context) Stack {
	return ctx.Value(keyStack).(Stack)
}

// GetFrame retrieves the current Frame from Context.
func GetFrame(ctx context.Context) (*Frame, error) {
	stack := GetStack(ctx)
	f, err := stack.Peek()
	if err != nil {
		return nil, errors.New("no frame")
	}
	return f.(*Frame), nil
}

// GetArgument retrieves the argument at the specified position from Context.
func GetArgument(ctx context.Context, idx int) (interface{}, error) {
	f, err := GetFrame(ctx)
	if err != nil {
		return nil, err
	}

	if idx >= len(f.Arguments) {
		return nil, errors.New("argument out of range")
	}

	return f.Arguments[idx], nil
}

// GetLocals retrieves the local variables from Context.
func GetLocals(ctx context.Context) (Heap, error) {
	f, err := GetFrame(ctx)
	if err != nil {
		return nil, err
	}
	return f.Locals, nil
}

func getResult(ctx context.Context) interface{} {
	return *ctx.Value(keyResult).(*interface{})
}

func setResult(ctx context.Context, res interface{}) {
	r := ctx.Value(keyResult).(*interface{})
	*r = res
}

type programContextKey int

const (
	keyLabels programContextKey = iota
)

func newProgramContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyLabels, map[string]int{})
	return ctx
}

// GetLabels retrieves the program labels from Context.
func GetLabels(ctx context.Context) map[string]int {
	return ctx.Value(keyLabels).(map[string]int)
}
