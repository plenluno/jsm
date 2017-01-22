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
	Load(k string) (interface{}, error)

	// Store stores the specified value under the specified key.
	Store(k string, v interface{})
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
	data, err := json.Marshal(h)
	return data, errors.Wrap(err, "failed to dump heap")
}

func (h *heap) Restore(data []byte) error {
	return errors.Wrap(json.Unmarshal(data, h), "failed to restore heap")
}
