package jsm

type frame struct {
	Arguments []interface{} `json:"arguments"`
	Locals    *heap         `json:"locals"`
	Operands  *stack        `json:"operands"`
	ReturnTo  *address      `json:"returnTo"`
}

func newFrame() *frame {
	f := new(frame)
	f.Locals = newHeap()
	f.Operands = newStack()
	f.ReturnTo = newAddress()
	return f
}
