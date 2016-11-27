package jsm

// Restorable objects can dump their entire states and
// restore themselves from the dump data.
type Restorable interface {
	Dump() ([]byte, error)
	Restore(data []byte) error
}
