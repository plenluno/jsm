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
		keyResult: new(interface{}),
	}
}

// GetPC retrieves the program counter.
func GetPC(ctx context.Context) Address {
	return (*ctx.(*machineContext))[keyPC].(*address)
}

// GetGlobalHeap retrieves the global heap.
func GetGlobalHeap(ctx context.Context) Heap {
	return (*ctx.(*machineContext))[keyHeap].(*heap)
}

func getCallStack(ctx context.Context) Stack {
	return (*ctx.(*machineContext))[keyStack].(*stack)
}

func getFrame(ctx context.Context) (*frame, error) {
	stack := getCallStack(ctx)
	f, err := stack.Peek()
	if err != nil {
		return nil, errors.New("no frame")
	}
	return f.(*frame), nil
}

// GetArgument retrieves the argument at the specified position.
func GetArgument(ctx context.Context, idx int) (interface{}, error) {
	f, err := getFrame(ctx)
	if err != nil {
		return nil, err
	}

	if idx >= len(f.Arguments) {
		return nil, errors.New("argument out of range")
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

func getResult(ctx context.Context) interface{} {
	return *(*ctx.(*machineContext))[keyResult].(*interface{})
}

func setResult(ctx context.Context, res interface{}) {
	r := (*ctx.(*machineContext))[keyResult].(*interface{})
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
