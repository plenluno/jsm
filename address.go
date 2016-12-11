package jsm

// Address is an address in JSM.
type Address interface {
	Clearable

	// GetValue returns the current address value.
	GetValue() int

	// SetValue sets the address to the given value.
	SetValue(addr int)

	// Increment increments the address value.
	Increment()
}

// NewAddress creates a new Address.
func NewAddress() Address {
	return newAddress()
}

type address int

func newAddress() *address {
	var addr address
	return &addr
}

func (a *address) GetValue() int {
	return int(*a)
}

func (a *address) SetValue(addr int) {
	*a = address(addr)
}

func (a *address) Increment() {
	*a++
}

func (a *address) Clear() {
	*a = 0
}
