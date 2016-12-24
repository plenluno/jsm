package jsm

import (
	"context"

	"github.com/pkg/errors"
)

// Process executes operations on JSM.
type Process func(ctx context.Context, immediates []interface{}) error

type processor map[Mnemonic]Process

func newProcessor() processor {
	p := processor{
		MnemonicPush:           push,
		MnemonicPop:            pop,
		MnemonicCall:           call,
		MnemonicReturn:         ret,
		MnemonicJump:           jmp,
		MnemonicJumpIfTrue:     jt,
		MnemonicJumpIfFalse:    jf,
		MnemonicEqual:          eq,
		MnemonicNotEqual:       ne,
		MnemonicGreaterThan:    gt,
		MnemonicGreaterOrEqual: ge,
		MnemonicLessThan:       lt,
		MnemonicLessOrEqual:    le,
	}
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

// Push pushes a value onto the operand stack.
func Push(ctx context.Context, v interface{}) error {
	frame, err := GetFrame(ctx)
	if err != nil {
		return err
	}

	frame.Operands.Push(v)
	return nil
}

// MultiPush pushes multiple values onto the operand stack.
func MultiPush(ctx context.Context, vs []interface{}) error {
	frame, err := GetFrame(ctx)
	if err != nil {
		return err
	}

	for _, v := range vs {
		frame.Operands.Push(v)
	}
	return nil
}

// Pop pops a value from the operand stack.
func Pop(ctx context.Context) (interface{}, error) {
	frame, err := GetFrame(ctx)
	if err != nil {
		return nil, err
	}

	v, err := frame.Operands.Pop()
	if err != nil {
		return nil, errors.New("no operand")
	}
	return v, nil
}

// MultiPop pops multiple values from the operand stack.
func MultiPop(ctx context.Context, n int) ([]interface{}, error) {
	frame, err := GetFrame(ctx)
	if err != nil {
		return nil, err
	}

	operands := make([]interface{}, n)
	for i := 0; i < n; i++ {
		v, err := frame.Operands.Pop()
		if err != nil {
			return nil, errors.New("no operand")
		}

		operands[n-i-1] = v
	}
	return operands, nil
}

func getAddress(vs []interface{}, idx int) (int, error) {
	if len(vs) <= idx {
		return -1, errors.New("no address")
	}

	addr, ok := vs[idx].(int)
	if !ok || addr < 0 {
		return -1, errors.New("invalid address")
	}

	return addr, nil
}

func getCount(vs []interface{}, idx, min int) (int, error) {
	if len(vs) <= idx {
		return min, nil
	}

	count, ok := vs[idx].(int)
	if !ok || count < min {
		return -1, errors.New("invalid count")
	}

	return count, nil
}

func push(ctx context.Context, immediates []interface{}) error {
	if err := MultiPush(ctx, immediates); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func pop(ctx context.Context, immediates []interface{}) error {
	n, err := getCount(immediates, 0, 1)
	if err != nil {
		return err
	}

	if _, err := MultiPop(ctx, n); err != nil {
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

	argc, err := getCount(immediates, 1, 0)
	if err != nil {
		return err
	}

	argv, err := MultiPop(ctx, argc)
	if err != nil {
		return nil
	}

	pc := GetPC(ctx)
	pc.Increment()

	frame := NewFrame()
	frame.Arguments = argv
	frame.ReturnTo.SetValue(pc.GetValue())
	GetStack(ctx).Push(frame)

	pc.SetValue(addr)
	return nil
}

func ret(ctx context.Context, immediates []interface{}) error {
	n, err := getCount(immediates, 0, 0)
	if err != nil {
		return err
	}

	res, err := MultiPop(ctx, n)
	if err != nil {
		return err
	}

	frame, err := GetFrame(ctx)
	if err != nil {
		return err
	}

	GetPC(ctx).SetValue(frame.ReturnTo.GetValue())
	_, err = GetStack(ctx).Pop()
	if err != nil {
		return err
	}

	if err := MultiPush(ctx, res); err != nil {
		setResult(ctx, res)
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

	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	if ToBool(v) {
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

	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	if !ToBool(v) {
		GetPC(ctx).SetValue(addr)
	} else {
		GetPC(ctx).Increment()
	}
	return nil
}

func eq(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, Equal(vs[0], vs[1])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func ne(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, !Equal(vs[0], vs[1])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func gt(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, Less(vs[1], vs[0])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func ge(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, Less(vs[1], vs[0]) || Equal(vs[1], vs[0])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func lt(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, Less(vs[0], vs[1])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func le(ctx context.Context, immediates []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	if err := Push(ctx, Less(vs[0], vs[1]) || Equal(vs[0], vs[1])); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}
