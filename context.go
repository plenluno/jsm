package jsm

import (
	"context"

	"github.com/pkg/errors"
)

type contextKey int

const (
	keyPC contextKey = iota
	keyHeap
	keyStack
	keyResult
)

func newContext(m *machine) context.Context {
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

func getResult(ctx context.Context) interface{} {
	return *ctx.Value(keyResult).(*interface{})
}

func setResult(ctx context.Context, res interface{}) {
	r := ctx.Value(keyResult).(*interface{})
	*r = res
}
