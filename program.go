package jsm

// Program is a program to run on JSM.
type Program struct {
	Instructions []Instruction `json:"program"`
}

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
	MnemonicPush   Mnemonic = "PUSH"
	MnemonicPop             = "POP"
	MnemonicCall            = "CALL"
	MnemonicReturn          = "RET"
)
