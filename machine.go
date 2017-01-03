package jsm

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

// Machine is a JSM.
type Machine interface {
	Clearable
	Restorable

	Run(program []Instruction, args []interface{}) (interface{}, error)

	Extend(mnemonic Mnemonic, process Process, preprocess Preprocess) error
}

// NewMachine creates a new Machine.
func NewMachine() Machine {
	return newMachine()
}

type machine struct {
	processor    *processor
	preprocessor *preprocessor

	Program []Instruction `json:"program"`
	PC      *address      `json:"pc"`
	Heap    *heap         `json:"heap"`
	Stack   *stack        `json:"stack"`

	context context.Context
}

func newMachine() *machine {
	m := new(machine)
	m.processor = newProcessor()
	m.preprocessor = newPreprocessor()
	m.PC = newAddress()
	m.Heap = newHeap()
	m.Stack = newStack()
	m.context = newMachineContext(m)
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
	p, err := m.preprocessor.preprocess(program)
	if err != nil {
		return err
	}

	if args == nil {
		args = []interface{}{}
	}

	m.Clear()
	m.Program = p

	frame := NewFrame()
	frame.Arguments = args
	frame.ReturnTo.SetValue(len(p))
	m.Stack.Push(frame)
	return nil
}

func (m *machine) inProgress() bool {
	pc := m.PC.GetValue()
	return pc >= 0 && pc < len(m.Program)
}

func (m *machine) step() error {
	return m.processor.process(m.context, &m.Program[m.PC.GetValue()])
}

func (m *machine) Extend(mnemonic Mnemonic, process Process, preprocess Preprocess) error {
	if err := m.processor.extend(mnemonic, process); err != nil {
		return err
	}
	return m.preprocessor.extend(mnemonic, preprocess)
}

func (m *machine) Clear() {
	m.Program = nil
	m.PC.Clear()
	m.Heap.Clear()
	m.Stack.Clear()
}

func (m *machine) Dump() ([]byte, error) {
	data, err := json.Marshal(m)
	return data, errors.Wrap(err, "failed to dump machine")
}

func (m *machine) Restore(data []byte) error {
	return errors.Wrap(json.Unmarshal(data, m), "failed to restore machine")
}
