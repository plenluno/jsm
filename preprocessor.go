package jsm

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

// Preprocess preprocesses the immediates of an instruction.
type Preprocess func(ctx context.Context, imms []interface{}) ([]interface{}, error)

type preprocessor map[Mnemonic]Preprocess

func newPreprocessor() *preprocessor {
	return &preprocessor{
		MnemonicPush:        noPreprocessing,
		MnemonicPop:         atMostOneInteger,
		MnemonicCall:        immediatesOfCall,
		MnemonicReturn:      atMostOneInteger,
		MnemonicJump:        oneAddress,
		MnemonicJumpIfTrue:  oneAddress,
		MnemonicJumpIfFalse: oneAddress,
		MnemonicAdd:         atMostOneNumber,
		MnemonicSubtract:    atMostOneNumber,
		MnemonicMultiply:    atMostOneNumber,
		MnemonicDivide:      atMostOneNumber,
	}
}

func (pp preprocessor) extend(mnemonic Mnemonic, preprocess Preprocess) error {
	if mnemonic == "" {
		return errors.New("no mnemonic")
	}

	if preprocess == nil {
		return nil
	}

	if _, ok := pp[mnemonic]; ok {
		return errors.Errorf("mnemonic already defined: %s", mnemonic)
	}

	pp[mnemonic] = preprocess
	return nil
}

func (pp preprocessor) preprocess(program []Instruction) ([]Instruction, error) {
	if program == nil {
		return nil, errors.New("no program")
	}

	ctx := newProgramContext()
	labels := GetLabels(ctx)
	for idx, inst := range program {
		if inst.Label != "" {
			labels[inst.Label] = idx
		}
	}

	preprocessed := make([]Instruction, len(program))
	for idx, inst := range program {
		m := inst.Mnemonic
		setMnemonic(ctx, m)

		p := pp[m]
		if p == nil {
			p = noImmediate
		}

		imms, err := p(ctx, inst.Immediates)
		if err != nil {
			return nil, err
		}

		preprocessed[idx] = Instruction{
			Mnemonic:   m,
			Immediates: imms,
			opcode:     opcode(m),
		}
	}
	return preprocessed, nil
}

func immediatesOfCall(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	switch len(imms) {
	case 0:
		return nil, preprocessingError(ctx, imms, "no immediate")
	case 1:
		return []interface{}{toAddress(ctx, imms[0])}, nil
	case 2:
		return []interface{}{toAddress(ctx, imms[0]), ToInteger(imms[1])}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func noPreprocessing(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	return imms, nil
}

func atMostOneInteger(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []interface{}{ToInteger(imms[0])}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func atMostOneNumber(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []interface{}{ToNumber(imms[0])}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func oneAddress(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	switch len(imms) {
	case 0:
		return nil, preprocessingError(ctx, imms, "no immediate")
	case 1:
		return []interface{}{toAddress(ctx, imms[0])}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func noImmediate(ctx context.Context, imms []interface{}) ([]interface{}, error) {
	if len(imms) > 0 {
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
	return nil, nil
}

func toAddress(ctx context.Context, v interface{}) interface{} {
	switch v.(type) {
	case string:
		return GetLabels(ctx)[v.(string)]
	default:
		return ToInteger(v)
	}
}

func preprocessingError(ctx context.Context, imms []interface{}, msg string) error {
	data, err := json.Marshal(Instruction{
		Mnemonic:   GetMnemonic(ctx),
		Immediates: imms,
	})
	if err != nil {
		return errors.Wrap(err, "cannot convert to json")
	}
	return errors.Errorf(msg+": %s", string(data))
}
