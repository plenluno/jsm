package jsm

// Instruction is a instruction for JSM.
type Instruction struct {
	Label      string        `json:"label,omitempty"`
	Mnemonic   Mnemonic      `json:"mnemonic"`
	Immediates []interface{} `json:"immediates,omitempty"`
	Comment    string        `json:"comment,omitempty"`
}

// Mnemonic is a instruction mnemonic for JSM.
type Mnemonic string

// These constants are instruction mnemonics.
const (
	MnemonicNop            Mnemonic = "nop"
	MnemonicPush                    = "push"
	MnemonicPop                     = "pop"
	MnemonicLoad                    = "ld"
	MnemonicLoadArgument            = "lda"
	MnemonicLoadLocal               = "ldl"
	MnemonicStore                   = "st"
	MnemonicStoreLocal              = "stl"
	MnemonicCall                    = "call"
	MnemonicReturn                  = "ret"
	MnemonicJump                    = "jmp"
	MnemonicJumpIfTrue              = "jt"
	MnemonicJumpIfFalse             = "jf"
	MnemonicEqual                   = "eq"
	MnemonicNotEqual                = "ne"
	MnemonicGreaterThan             = "gt"
	MnemonicGreaterOrEqual          = "ge"
	MnemonicLessThan                = "lt"
	MnemonicLessOrEqual             = "le"
	MnemonicNot                     = "not"
	MnemonicAnd                     = "and"
	MnemonicOr                      = "or"
	MnemonicAdd                     = "add"
	MnemonicSubtract                = "sub"
	MnemonicMultiply                = "mul"
	MnemonicDivide                  = "div"
)
