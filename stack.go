package jsm

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Stack is a stack of Value’s.
type Stack interface {
	Clearable
	Restorable

	// Push pushes a value onto the stack.
	Push(v Value)

	// MultiPush pushes multiple values onto the stack.
	MultiPush(vs []Value)

	// Pop pops a value from the stack.
	Pop() (Value, error)

	// MultiPop pops multiple values from the stack.
	// The return slice can be modified by Push/MultiPush.
	// So you should read the slice before calling Push/MultiPush.
	MultiPop(n int) ([]Value, error)

	// Peek returns the value at the top of the stack.
	Peek() (Value, error)

	// Do executes the given operation against the values at the top of the stack.
	Do(op func([]Value) (Value, error), arity int) error
}

type stack []Value

func newStack() *stack {
	s := make(stack, 0, 10)
	return &s
}

func (s *stack) Push(v Value) {
	*s = append(*s, v)
}

func (s *stack) MultiPush(vs []Value) {
	*s = append(*s, vs...)
}

func (s *stack) Pop() (Value, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	v := (*s)[l-1]
	*s = (*s)[:l-1]
	return v, nil
}

func (s *stack) MultiPop(n int) ([]Value, error) {
	l := len(*s)
	if l < n {
		return nil, errors.New("too few elements")
	}

	vs := (*s)[l-n:]
	*s = (*s)[:l-n]
	return vs, nil
}

func (s *stack) Peek() (Value, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	return (*s)[l-1], nil
}

func (s *stack) Do(op func([]Value) (Value, error), arity int) error {
	l := len(*s)
	if l < arity {
		return errors.New("too few elements")
	}

	v, err := op((*s)[l-arity:])
	if err != nil {
		return err
	}

	(*s)[l-arity] = v
	*s = (*s)[:l-arity+1]
	return nil
}

func (s *stack) Clear() {
	*s = (*s)[:0]
}

func (s *stack) Dump() ([]byte, error) {
	data, err := json.Marshal(s)
	return data, errors.Wrap(err, "failed to dump stack")
}

func (s *stack) Restore(data []byte) error {
	return errors.Wrap(json.Unmarshal(data, s), "failed to restore stack")
}
