package jsm

// Frame is a frame in JSM call stack.
type Frame struct {
	Locals   *heap  `json:"locals"`
	Operands *stack `json:"operands"`
	Return   Return `json:"return"`
}

// Return holds a return address and a return value.
type Return struct {
	Address *address    `json:"address"`
	Value   interface{} `json:"value,omitempty"`
}

// NewFrame creates a new Frame.
func NewFrame() *Frame {
	f := new(Frame)
	f.Locals = newHeap()
	f.Operands = newStack()
	f.Return.Address = newAddress()
	return f
}
