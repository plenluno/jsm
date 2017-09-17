package jsm

// Type represents the type of a JSON value.
type Type int

// These constants are the types of JSON values.
const (
	TypeUndefined Type = iota
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeArray
	TypeObject

	// TypePointer is a special type that represents a pointer to an arbitrary type,
	// but is always marshaled as 'null' in JSON.
	TypePointer
)
