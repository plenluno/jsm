package jsm

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
)

// Process executes an operation on JSM.
type Process func(ctx context.Context, immediates []interface{}) error

type processor map[Mnemonic]Process

func newProcessor() processor {
	p := processor{}
	p[MnemonicPush] = push
	p[MnemonicPop] = pop
	p[MnemonicCall] = call
	p[MnemonicReturn] = ret
	p[MnemonicJump] = jmp
	p[MnemonicJumpIfTrue] = jt
	p[MnemonicJumpIfFalse] = jf
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
		return errors.Errorf("mnemonic already defined: %s", mnemonic)
	}

	p[mnemonic] = process
	return nil
}

func isTrue(v interface{}) bool {
	return !isFalse(v)
}

func isFalse(v interface{}) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Bool:
		return val.Bool() == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0.0
	case reflect.String:
		return val.String() == ""
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

func popOperands(frame *Frame, n int) ([]interface{}, error) {
	operands := make([]interface{}, n)
	for i := 0; i < n; i++ {
		v, err := frame.Operands.Pop()
		if err != nil {
			return nil, err
		}

		operands[i] = v
	}
	return operands, nil
}

func getAddress(vs []interface{}, idx int) (int, error) {
	if len(vs) > idx {
		return -1, errors.New("no address")
	}

	addr, ok := vs[idx].(int)
	if !ok {
		return -1, errors.New("invalid address")
	}

	return addr, nil
}

func push(ctx context.Context, immediates []interface{}) error {
	frame := GetFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	for _, imm := range immediates {
		frame.Operands.Push(imm)
	}

	GetPC(ctx).Increment()
	return nil
}

func pop(ctx context.Context, immediates []interface{}) error {
	n := 1
	if len(immediates) > 0 {
		var ok bool
		n, ok = immediates[0].(int)
		if !ok || n < 1 {
			return errors.New("invalid pop count")
		}
	}

	frame := GetFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	if _, err := popOperands(frame, n); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func call(ctx context.Context, immediates []interface{}) error {
	addr, err := getAddress(immediates, 0)
	if err != nil {
		return err
	}

	argc := 0
	if len(immediates) > 1 {
		var ok bool
		argc, ok = immediates[1].(int)
		if !ok || argc < 0 {
			return errors.New("invalid argument count")
		}
	}

	frame := GetFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	argv, err := popOperands(frame, argc)
	if err != nil {
		return nil
	}

	pc := GetPC(ctx)
	pc.Increment()

	frame = NewFrame()
	frame.Arguments = argv
	frame.ReturnTo.SetValue(pc.GetValue())
	GetStack(ctx).Push(frame)

	pc.SetValue(addr)
	return nil
}

func ret(ctx context.Context, immediates []interface{}) error {
	n := 0
	if len(immediates) > 0 {
		var ok bool
		n, ok = immediates[0].(int)
		if !ok || n < 0 {
			return errors.New("invalid return value count")
		}
	}

	v, err := GetStack(ctx).Pop()
	if err != nil {
		return err
	}

	frame, ok := v.(*Frame)
	if !ok {
		return errors.New("invalid frame")
	}

	GetPC(ctx).SetValue(frame.ReturnTo.GetValue())

	retVals, err := popOperands(frame, n)
	if err != nil {
		return err
	}

	frame = GetFrame(ctx)
	if frame == nil {
		setResult(ctx, retVals)
		return nil
	}

	for _, v := range retVals {
		frame.Operands.Push(v)
	}
	return nil
}

func jmp(ctx context.Context, immediates []interface{}) error {
	addr, err := getAddress(immediates, 0)
	if err != nil {
		return err
	}

	GetPC(ctx).SetValue(addr)
	return nil
}

func jt(ctx context.Context, immediates []interface{}) error {
	addr, err := getAddress(immediates, 0)
	if err != nil {
		return err
	}

	frame := GetFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	v, err := frame.Operands.Pop()
	if err != nil {
		return err
	}

	if isTrue(v) {
		GetPC(ctx).SetValue(addr)
	} else {
		GetPC(ctx).Increment()
	}
	return nil
}

func jf(ctx context.Context, immediates []interface{}) error {
	addr, err := getAddress(immediates, 0)
	if err != nil {
		return err
	}

	frame := GetFrame(ctx)
	if frame == nil {
		return errors.New("no frame")
	}

	v, err := frame.Operands.Pop()
	if err != nil {
		return err
	}

	if isFalse(v) {
		GetPC(ctx).SetValue(addr)
	} else {
		GetPC(ctx).Increment()
	}
	return nil
}
