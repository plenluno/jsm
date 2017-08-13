package jsm

import "github.com/pkg/errors"

type frame struct {
	Arguments []Value `json:"arguments"`
	Locals    *heap   `json:"locals"`
	Operands  *stack  `json:"operands"`
	ReturnTo  int     `json:"returnTo"`
}

func newFrame() *frame {
	f := new(frame)
	f.Locals = newHeap()
	f.Operands = newStack()
	return f
}

type callStack []*frame

func newCallStack() *callStack {
	cs := make(callStack, 0, 10)
	return &cs
}

func (cs *callStack) Push(f *frame) {
	*cs = append(*cs, f)
}

func (cs *callStack) Pop() (*frame, error) {
	l := len(*cs)
	if l == 0 {
		return nil, errors.New("no frame")
	}

	f := (*cs)[l-1]
	*cs = (*cs)[:l-1]
	return f, nil
}

func (cs *callStack) Peek() (*frame, error) {
	l := len(*cs)
	if l == 0 {
		return nil, errors.New("no frame")
	}
	return (*cs)[l-1], nil
}

func (cs *callStack) Clear() {
	*cs = (*cs)[:0]
}
