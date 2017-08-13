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
)
