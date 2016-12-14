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
	MnemonicPush   Mnemonic = "push"
	MnemonicPop             = "pop"
	MnemonicCall            = "call"
	MnemonicReturn          = "ret"
)
