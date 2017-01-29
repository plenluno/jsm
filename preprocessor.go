package jsm

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

// Preprocess preprocesses the immediates of an instruction.
type Preprocess func(ctx context.Context, imms []Value) ([]Value, error)

type preprocessor map[Mnemonic]Preprocess

func newPreprocessor() *preprocessor {
	return &preprocessor{
		MnemonicPush:           noPreprocessing,
		MnemonicPop:            atMostOneInteger,
		MnemonicLoad:           atMostOneString,
		MnemonicLoadArgument:   atMostOneInteger,
		MnemonicLoadLocal:      atMostOneString,
		MnemonicStore:          immediatesOfStore,
		MnemonicStoreLocal:     immediatesOfStore,
		MnemonicCall:           immediatesOfCall,
		MnemonicReturn:         atMostOneInteger,
		MnemonicJump:           oneAddress,
		MnemonicJumpIfTrue:     oneAddress,
		MnemonicJumpIfFalse:    oneAddress,
		MnemonicEqual:          atMostOneImmediate,
		MnemonicNotEqual:       atMostOneImmediate,
		MnemonicGreaterThan:    atMostOneImmediate,
		MnemonicGreaterOrEqual: atMostOneImmediate,
		MnemonicLessThan:       atMostOneImmediate,
		MnemonicLessOrEqual:    atMostOneImmediate,
		MnemonicAdd:            atMostOneNumber,
		MnemonicSubtract:       atMostOneNumber,
		MnemonicMultiply:       atMostOneNumber,
		MnemonicDivide:         atMostOneNumber,
		MnemonicIncrement:      atMostOneString,
		MnemonicIncrementLocal: atMostOneString,
		MnemonicDecrement:      atMostOneString,
		MnemonicDecrementLocal: atMostOneString,
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

func noPreprocessing(ctx context.Context, imms []Value) ([]Value, error) {
	return imms, nil
}

func noImmediate(ctx context.Context, imms []Value) ([]Value, error) {
	if len(imms) > 0 {
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
	return nil, nil
}

func atMostOneImmediate(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return imms, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func atMostOneInteger(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []Value{IntegerValue(ToInteger(imms[0]))}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func atMostOneNumber(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []Value{NumberValue(ToNumber(imms[0]))}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func atMostOneString(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []Value{StringValue(ToString(imms[0]))}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func oneAddress(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, preprocessingError(ctx, imms, "no immediate")
	case 1:
		return []Value{toAddress(ctx, imms[0])}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func immediatesOfCall(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, preprocessingError(ctx, imms, "no immediate")
	case 1:
		return []Value{toAddress(ctx, imms[0])}, nil
	case 2:
		return []Value{toAddress(ctx, imms[0]), IntegerValue(ToInteger(imms[1]))}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func immediatesOfStore(ctx context.Context, imms []Value) ([]Value, error) {
	switch len(imms) {
	case 0:
		return nil, nil
	case 1:
		return []Value{StringValue(ToString(imms[0]))}, nil
	case 2:
		return []Value{StringValue(ToString(imms[0])), imms[1]}, nil
	default:
		return nil, preprocessingError(ctx, imms, "too many immediates")
	}
}

func toAddress(ctx context.Context, v Value) Value {
	switch v.(type) {
	case string:
		return IntegerValue(GetLabels(ctx)[ToString(v)])
	default:
		return IntegerValue(ToInteger(v))
	}
}

func preprocessingError(ctx context.Context, imms []Value, msg string) error {
	data, err := json.Marshal(Instruction{
		Mnemonic:   GetMnemonic(ctx),
		Immediates: imms,
	})
	if err != nil {
		return errors.Wrap(err, "cannot convert to json")
	}
	return errors.Errorf(msg+": %s", string(data))
}
