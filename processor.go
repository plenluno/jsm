package jsm

import (
	"context"
	"errors"
	"fmt"
)

type Process func(ctx context.Context, operands []interface{}) error

type processor map[Mnemonic]Process

func newProcessor() processor {
	p := processor{}
	p[MnemonicPush] = push
	p[MnemonicPop] = pop
	return p
}

func (p processor) extend(mnemonic Mnemonic, process Process) error {
	if mnemonic == "" {
		return errors.New("no mnemonic")
	}

	if process == nil {
		return errors.New("no process")
	}

	if _, ok := p[mnemonic]; ok {
		return fmt.Errorf("mnemonic already defined: %s", mnemonic)
	}

	p[mnemonic] = process
	return nil
}

func incrementPC(ctx context.Context) error {
	pc := ExtractPC(ctx)
	pc.Increment()
	return nil
}

func push(ctx context.Context, operands []interface{}) error {
	frame := ExtractFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	for _, o := range operands {
		frame.Operands.Push(o)
	}

	return incrementPC(ctx)
}

func pop(ctx context.Context, operands []interface{}) error {
	frame := ExtractFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	if _, err := frame.Operands.Pop(); err != nil {
		return err
	}

	return incrementPC(ctx)
}
