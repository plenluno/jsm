package jsm

// Preprocess preprocesses the given program.
func Preprocess(program []Instruction) ([]Instruction, error) {
	// TODO: Inspect program
	return preprocess(program), nil
}

func preprocess(program []Instruction) []Instruction {
	preprocessed := make([]Instruction, len(program))
	addrs := map[string]int{}
	for idx, inst := range program {
		preprocessed[idx].Mnemonic = inst.Mnemonic
		if inst.Label != "" {
			addrs[inst.Label] = idx
		}
	}
	for idx, inst := range program {
		preprocessed[idx].Immediates = preprocessImmediates(inst, addrs)
	}
	return preprocessed
}

func preprocessImmediates(inst Instruction, addrs map[string]int) []interface{} {
	preprocessed := make([]interface{}, len(inst.Immediates))
	for idx, imm := range inst.Immediates {
		switch inst.Mnemonic {
		case MnemonicPop, MnemonicReturn:
			preprocessed[idx] = ToInteger(imm)
		case MnemonicCall:
			if idx == 0 {
				preprocessed[idx] = preprocessAddress(imm, addrs)
			} else {
				preprocessed[idx] = ToInteger(imm)
			}

		case MnemonicJump, MnemonicJumpIfTrue, MnemonicJumpIfFalse:
			preprocessed[idx] = preprocessAddress(imm, addrs)
		default:
			preprocessed[idx] = imm
		}
	}
	return preprocessed
}

func preprocessAddress(v interface{}, addrs map[string]int) interface{} {
	switch v.(type) {
	case string:
		return addrs[v.(string)]
	default:
		return ToInteger(v)
	}
}
