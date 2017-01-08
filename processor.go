package jsm

import (
	"context"

	"github.com/pkg/errors"
)

// Process executes operations on JSM.
type Process func(ctx context.Context, imms []interface{}) error

type processor []Process

func newProcessor() *processor {
	p := new(processor)
	extend := func(mnemonic Mnemonic, process Process) {
		if err := p.extend(mnemonic, process); err != nil {
			panic(err)
		}
	}
	extend(MnemonicNop, nop)
	extend(MnemonicPush, push)
	extend(MnemonicPop, pop)
	extend(MnemonicLoad, ld)
	extend(MnemonicLoadArgument, lda)
	extend(MnemonicLoadLocal, ldl)
	extend(MnemonicStore, st)
	extend(MnemonicStoreLocal, stl)
	extend(MnemonicCall, call)
	extend(MnemonicReturn, ret)
	extend(MnemonicJump, jmp)
	extend(MnemonicJumpIfTrue, jt)
	extend(MnemonicJumpIfFalse, jf)
	extend(MnemonicEqual, binaryOp(eq))
	extend(MnemonicNotEqual, binaryOp(ne))
	extend(MnemonicGreaterThan, binaryOp(gt))
	extend(MnemonicGreaterOrEqual, binaryOp(ge))
	extend(MnemonicLessThan, binaryOp(lt))
	extend(MnemonicLessOrEqual, binaryOp(le))
	extend(MnemonicNot, unaryOp(not))
	extend(MnemonicAnd, binaryOp(and))
	extend(MnemonicOr, binaryOp(or))
	extend(MnemonicNeg, unaryOp(neg))
	extend(MnemonicAdd, binaryOp(add))
	extend(MnemonicSubtract, binaryOp(sub))
	extend(MnemonicMultiply, binaryOp(mul))
	extend(MnemonicDivide, binaryOp(div))
	return p
}

func (p *processor) extend(mnemonic Mnemonic, process Process) error {
	if mnemonic == "" {
		return errors.New("no mnemonic")
	}

	if process == nil {
		return errors.New("no process")
	}

	oc := opcode(mnemonic)
	size := len(*p)
	if size > oc && (*p)[oc] != nil {
		return errors.Errorf("%s already defined", mnemonic)
	}

	if size <= oc {
		tmp := *p
		*p = make([]Process, oc+1)
		copy(*p, tmp)
	}
	(*p)[oc] = process
	return nil
}

func (p processor) process(ctx context.Context, inst *Instruction) error {
	oc := inst.opcode
	if len(p) <= oc || p[oc] == nil {
		return errors.Errorf("cannot process %s", inst.Mnemonic)
	}
	return p[oc](ctx, inst.Immediates)
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

	frame.Operands.MultiPush(vs)
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

	operands, err := frame.Operands.MultiPop(n)
	if err != nil {
		return nil, errors.New("too few operands")
	}
	return operands, nil
}

// Do executes the given operation against the values at the top of the operand stack.
func Do(ctx context.Context, op func([]interface{}) (interface{}, error), arity int) error {
	frame, err := GetFrame(ctx)
	if err != nil {
		return err
	}

	if err := frame.Operands.Do(op, arity); err != nil {
		return errors.New("too few operands")
	}
	return nil
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

func nop(ctx context.Context, imms []interface{}) error {
	GetPC(ctx).Increment()
	return nil
}

func push(ctx context.Context, imms []interface{}) error {
	if err := MultiPush(ctx, imms); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func pop(ctx context.Context, imms []interface{}) error {
	n, err := getCount(imms, 0, 1)
	if err != nil {
		return err
	}

	if _, err := MultiPop(ctx, n); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func ld(ctx context.Context, imms []interface{}) error {
	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	v, _ = GetHeap(ctx).Load(ToString(v))
	if err := Push(ctx, v); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func lda(ctx context.Context, imms []interface{}) error {
	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	a, err := GetArgument(ctx, ToInteger(v))
	if err != nil {
		return err
	}

	if err := Push(ctx, a); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func ldl(ctx context.Context, imms []interface{}) error {
	v, err := Pop(ctx)
	if err != nil {
		return err
	}

	ls, err := GetLocals(ctx)
	if err != nil {
		return err
	}

	v, _ = ls.Load(ToString(v))
	if err := Push(ctx, v); err != nil {
		return err
	}

	GetPC(ctx).Increment()
	return nil
}

func st(ctx context.Context, imms []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	GetHeap(ctx).Store(ToString(vs[0]), vs[1])
	GetPC(ctx).Increment()
	return nil
}

func stl(ctx context.Context, imms []interface{}) error {
	vs, err := MultiPop(ctx, 2)
	if err != nil {
		return err
	}

	ls, err := GetLocals(ctx)
	if err != nil {
		return err
	}

	ls.Store(ToString(vs[0]), vs[1])
	GetPC(ctx).Increment()
	return nil
}

func call(ctx context.Context, imms []interface{}) error {
	addr, err := getAddress(imms, 0)
	if err != nil {
		return err
	}

	argc, err := getCount(imms, 1, 0)
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

func ret(ctx context.Context, imms []interface{}) error {
	n, err := getCount(imms, 0, 0)
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

func jmp(ctx context.Context, imms []interface{}) error {
	addr, err := getAddress(imms, 0)
	if err != nil {
		return err
	}

	GetPC(ctx).SetValue(addr)
	return nil
}

func jt(ctx context.Context, imms []interface{}) error {
	addr, err := getAddress(imms, 0)
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

func jf(ctx context.Context, imms []interface{}) error {
	addr, err := getAddress(imms, 0)
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

func unaryOp(op func([]interface{}) (interface{}, error)) Process {
	return func(ctx context.Context, imms []interface{}) error {
		if err := Do(ctx, op, 1); err != nil {
			return err
		}

		GetPC(ctx).Increment()
		return nil
	}
}

func binaryOp(op func([]interface{}) (interface{}, error)) Process {
	return func(ctx context.Context, imms []interface{}) error {
		if len(imms) > 0 {
			if err := Push(ctx, imms[0]); err != nil {
				return err
			}
		}

		if err := Do(ctx, op, 2); err != nil {
			return err
		}

		GetPC(ctx).Increment()
		return nil
	}
}

func eq(vs []interface{}) (interface{}, error) {
	return Equal(vs[0], vs[1]), nil
}

func ne(vs []interface{}) (interface{}, error) {
	return !Equal(vs[0], vs[1]), nil
}

func gt(vs []interface{}) (interface{}, error) {
	return Less(vs[1], vs[0]), nil
}

func ge(vs []interface{}) (interface{}, error) {
	return Less(vs[1], vs[0]) || Equal(vs[1], vs[0]), nil
}

func lt(vs []interface{}) (interface{}, error) {
	return Less(vs[0], vs[1]), nil
}

func le(vs []interface{}) (interface{}, error) {
	return Less(vs[0], vs[1]) || Equal(vs[0], vs[1]), nil
}

func not(vs []interface{}) (interface{}, error) {
	return !ToBool(vs[0]), nil
}

func and(vs []interface{}) (interface{}, error) {
	return ToBool(vs[0]) && ToBool(vs[1]), nil
}

func or(vs []interface{}) (interface{}, error) {
	return ToBool(vs[0]) || ToBool(vs[1]), nil
}

func neg(vs []interface{}) (interface{}, error) {
	return -ToNumber(vs[0]), nil
}

func add(vs []interface{}) (interface{}, error) {
	return ToNumber(vs[0]) + ToNumber(vs[1]), nil
}

func sub(vs []interface{}) (interface{}, error) {
	return ToNumber(vs[0]) - ToNumber(vs[1]), nil
}

func mul(vs []interface{}) (interface{}, error) {
	return ToNumber(vs[0]) * ToNumber(vs[1]), nil
}

func div(vs []interface{}) (interface{}, error) {
	num1 := ToNumber(vs[0])
	num2 := ToNumber(vs[1])
	if num2 == 0.0 {
		return nil, errors.New("divide by zero")
	}

	return num1 / num2, nil
}
