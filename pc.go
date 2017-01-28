package jsm

// ProgramCounter is the program counter of JSM.
type ProgramCounter interface {
	Clearable

	// Index returns the index of the current instruction.
	Index() int

	// SetIndex sets the program counter to the given index.
	SetIndex(idx int)

	// Increment increments the program counter.
	Increment()
}

type programCounter int

func newProgramCounter() *programCounter {
	var pc programCounter
	return &pc
}

func (pc *programCounter) Index() int {
	return int(*pc)
}

func (pc *programCounter) SetIndex(idx int) {
	*pc = programCounter(idx)
}

func (pc *programCounter) Increment() {
	*pc++
}

func (pc *programCounter) Clear() {
	*pc = 0
}
