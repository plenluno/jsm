package jsm

// Frame is a frame in JSM call stack.
type Frame struct {
	Arguments []interface{} `json:"arguments"`
	Locals    *heap         `json:"locals"`
	Operands  *stack        `json:"operands"`
	ReturnTo  *address      `json:"returnTo"`
}

// NewFrame creates a new Frame.
func NewFrame() *Frame {
	f := new(Frame)
	f.Locals = newHeap()
	f.Operands = newStack()
	f.ReturnTo = newAddress()
	return f
}
