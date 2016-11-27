package jsm

import "context"

type contextKey int

const (
	keyPC contextKey = iota
	keyHeap
	keyStack
)

func newContext(m *machine) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyPC, Address(m.pc))
	ctx = context.WithValue(ctx, keyHeap, Heap(m.heap))
	ctx = context.WithValue(ctx, keyStack, Stack(m.stack))
	return ctx
}

// ExtractPC extracts the program counter from Context.
func ExtractPC(ctx context.Context) Address {
	return ctx.Value(keyPC).(Address)
}

// ExtractHeap extracts Heap from Context.
func ExtractHeap(ctx context.Context) Heap {
	return ctx.Value(keyHeap).(Heap)
}

// ExtractStack extracts Stack from Context.
func ExtractStack(ctx context.Context) Stack {
	return ctx.Value(keyStack).(Stack)
}

// ExtractFrame extracts the current Frame from Context.
func ExtractFrame(ctx context.Context) *Frame {
	stack := ExtractStack(ctx)
	f, err := stack.Peek()
	if err != nil {
		return nil
	}
	return f.(*Frame)
}
