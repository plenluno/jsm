package jsm

import (
	"context"
	"encoding/json"
	"errors"
)

// Machine is a JSM.
type Machine interface {
	Clearable
	Restorable

	Run(program []Instruction, args []interface{}) (interface{}, error)

	Extend(mnemonic Mnemonic, process Process) error
}

// NewMachine creates a new Machine.
func NewMachine() Machine {
	return newMachine()
}

type machine struct {
	processor processor

	Program []Instruction `json:"program"`
	PC      *address      `json:"pc"`
	Heap    *heap         `json:"heap"`
	Stack   *stack        `json:"stack"`

	context context.Context
}

func newMachine() *machine {
	m := new(machine)
	m.processor = newProcessor()
	m.PC = newAddress()
	m.Heap = newHeap()
	m.Stack = newStack()
	m.context = newContext(m)
	return m
}

func (m *machine) Run(program []Instruction, args []interface{}) (interface{}, error) {
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

func (m *machine) load(program []Instruction, args []interface{}) error {
	if program == nil {
		return errors.New("no program")
	}

	if args == nil {
		args = []interface{}{}
	}

	// TODO: Inspect program

	m.Clear()
	m.Program = program

	frame := NewFrame()
	frame.Arguments = args
	frame.ReturnTo.SetValue(len(program))
	m.Stack.Push(frame)
	return nil
}

func (m *machine) inProgress() bool {
	pc := m.PC.GetValue()
	return pc >= 0 && pc < len(m.Program)
}

func (m *machine) step() error {
	inst := &m.Program[m.PC.GetValue()]
	return m.processor[inst.Mnemonic](m.context, inst.Immediates)
}

func (m *machine) Extend(mnemonic Mnemonic, process Process) error {
	return m.processor.extend(mnemonic, process)
}

func (m *machine) Clear() {
	m.Program = nil
	m.PC.Clear()
	m.Heap.Clear()
	m.Stack.Clear()
}

func (m *machine) Dump() ([]byte, error) {
	return json.Marshal(m)
}

func (m *machine) Restore(data []byte) error {
	return json.Unmarshal(data, m)
}
