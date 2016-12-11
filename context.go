package jsm

import "context"

type contextKey int

const (
	keyPC contextKey = iota
	keyHeap
	keyStack
	keyResult
)

func newContext(m *machine) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyPC, Address(m.pc))
	ctx = context.WithValue(ctx, keyHeap, Heap(m.heap))
	ctx = context.WithValue(ctx, keyStack, Stack(m.stack))
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
func GetFrame(ctx context.Context) *Frame {
	stack := GetStack(ctx)
	f, err := stack.Peek()
	if err != nil {
		return nil
	}
	return f.(*Frame)
}

func getResult(ctx context.Context) interface{} {
	return *ctx.Value(keyResult).(*interface{})
}

func setResult(ctx context.Context, res interface{}) {
	r := ctx.Value(keyResult).(*interface{})
	*r = res
}
