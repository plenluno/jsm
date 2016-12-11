package jsm

import (
	"context"
	"errors"
)

// Machine is a JSM.
type Machine interface {
	Clearable

	Run(program *Program, args []interface{}) (interface{}, error)

	Extend(mnemonic Mnemonic, process Process) error
}

// NewMachine creates a new Machine.
func NewMachine() Machine {
	return newMachine()
}

type machine struct {
	processor processor

	program *Program `json:"program"`
	pc      *address `json:"pc"`

	heap  *heap  `json:"heap"`
	stack *stack `json:"stack"`

	context context.Context
}

func newMachine() *machine {
	m := new(machine)
	m.processor = newProcessor()
	m.pc = newAddress()
	m.heap = newHeap()
	m.stack = newStack()
	m.context = newContext(m)
	return m
}

func (m *machine) Run(program *Program, args []interface{}) (interface{}, error) {
	if err := m.load(program, args); err != nil {
		return nil, err
	}

	for m.inProgress() {
		if err := m.step(); err != nil {
			return nil, err
		}
	}

	return getResult(m.context), nil
}

func (m *machine) load(program *Program, args []interface{}) error {
	if program == nil {
		return errors.New("no program")
	}

	if args == nil {
		args = []interface{}{}
	}

	// TODO: Inspect program

	m.Clear()
	m.program = program

	frame := NewFrame()
	frame.Arguments = args
	frame.ReturnTo.SetValue(len(program.Instructions))
	m.stack.Push(frame)
	return nil
}

func (m *machine) inProgress() bool {
	pc := m.pc.GetValue()
	return pc >= 0 && pc < len(m.program.Instructions)
}

func (m *machine) step() error {
	inst := &m.program.Instructions[m.pc.GetValue()]
	return m.processor[inst.Mnemonic](m.context, inst.Immediates)
}

func (m *machine) Extend(mnemonic Mnemonic, process Process) error {
	return m.processor.extend(mnemonic, process)
}

func (m *machine) Clear() {
	m.program = nil
	m.pc.Clear()
	m.heap.Clear()
	m.stack.Clear()
}
