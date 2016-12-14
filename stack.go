package jsm

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Stack is the stack of a JSM.
type Stack interface {
	Clearable
	Restorable

	Push(v interface{})
	Pop() (interface{}, error)
	Peek() (interface{}, error)
}

// NewStack creates a new Stack.
func NewStack() Stack {
	return newStack()
}

type stack []interface{}

func newStack() *stack {
	return &stack{}
}

func (s *stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *stack) Pop() (interface{}, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	v := (*s)[l-1]
	*s = (*s)[:l-1]
	return v, nil
}

func (s *stack) Peek() (interface{}, error) {
	l := len(*s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	return (*s)[l-1], nil
}

func (s *stack) Clear() {
	*s = []interface{}{}
}

func (s *stack) Dump() ([]byte, error) {
	data, err := json.Marshal(s)
	return data, errors.Wrap(err, "failed to dump stack")
}

func (s *stack) Restore(data []byte) error {
	return errors.Wrap(json.Unmarshal(data, s), "failed to restore stack")
}
