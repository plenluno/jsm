package jsm

// Address is an address in JSM.
type Address interface {
	Clearable

	// Value returns the current address value.
	Value() int

	// Increment increments the address.
	Increment()

	// Jump adds the given difference to the address.
	Jump(diff int)
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

func (a *address) Value() int {
	return int(*a)
}

func (a *address) Increment() {
	*a++
}

func (a *address) Jump(diff int) {
	*a += address(diff)
}

func (a *address) Clear() {
	*a = 0
}
