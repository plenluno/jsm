package jsm

import (
	"encoding/json"
	"errors"
)

// Heap is the heap of a JSM.
type Heap interface {
	Clearable
	Restorable

	Load(k string) (interface{}, error)
	Store(k string, v interface{})
}

// NewHeap creates a new Heap.
func NewHeap() Heap {
	return newHeap()
}

type heap map[string]interface{}

func newHeap() *heap {
	return &heap{}
}

func (h *heap) Load(k string) (interface{}, error) {
	v, ok := (*h)[k]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (h *heap) Store(k string, v interface{}) {
	(*h)[k] = v
}

func (h *heap) Clear() {
	*h = map[string]interface{}{}
}

func (h *heap) Dump() ([]byte, error) {
	return json.Marshal(h)
}

func (h *heap) Restore(data []byte) error {
	return json.Unmarshal(data, h)
}
