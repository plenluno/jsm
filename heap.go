package jsm

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Heap is a heap in JSM.
type Heap interface {
	Clearable
	Restorable

	// Load returns the value associated with the specified key,
	// or an error if the key has no associated value.
	Load(k string) (Value, error)

	// Store stores the specified value under the specified key.
	Store(k string, v Value)
}

type heap map[string]Value

func newHeap() *heap {
	return &heap{}
}

func (h *heap) Load(k string) (Value, error) {
	v, ok := (*h)[k]
	if !ok {
		return NullValue(), errors.New("not found")
	}
	return v, nil
}

func (h *heap) Store(k string, v Value) {
	(*h)[k] = v
}

func (h *heap) Clear() {
	*h = map[string]Value{}
}

func (h *heap) Dump() ([]byte, error) {
	data, err := json.Marshal(h)
	return data, errors.Wrap(err, "failed to dump heap")
}

func (h *heap) Restore(data []byte) error {
	return errors.Wrap(json.Unmarshal(data, h), "failed to restore heap")
}
