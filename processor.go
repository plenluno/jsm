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
	extend(MnemonicLoad, loadOp(ld))
	extend(MnemonicLoadArgument, loadOp(lda))
	extend(MnemonicLoadLocal, loadOp(ldl))
	extend(MnemonicStore, storeOp(st))
	extend(MnemonicStoreLocal, storeOp(stl))
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
	extend(MnemonicIncrement, loadStoreOp(inc))
	extend(MnemonicIncrementLocal, loadStoreOp(incl))
	extend(MnemonicDecrement, loadStoreOp(dec))
	extend(MnemonicDecrementLocal, loadStoreOp(decl))
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
	stack, err := GetOperandStack(ctx)
	if err != nil {
		return err
	}

	stack.Push(v)
	return nil
}

// MultiPush pushes multiple values onto the operand stack.
func MultiPush(ctx context.Context, vs []interface{}) error {
	stack, err := GetOperandStack(ctx)
	if err != nil {
		return err
	}

	stack.MultiPush(vs)
	return nil
}

// Pop pops a value from the operand stack.
func Pop(ctx context.Context) (interface{}, error) {
	stack, err := GetOperandStack(ctx)
	if err != nil {
		return nil, err
	}

	v, err := stack.Pop()
	if err != nil {
		return nil, errors.New("no operand")
	}
	return v, nil
}

// MultiPop pops multiple values from the operand stack.
func MultiPop(ctx context.Context, n int) ([]interface{}, error) {
	stack, err := GetOperandStack(ctx)
	if err != nil {
		return nil, err
	}

	operands, err := stack.MultiPop(n)
	if err != nil {
		return nil, errors.New("too few operands")
	}
	return operands, nil
}

// Do executes the given operation against the values at the top of the operand stack.
func Do(ctx context.Context, op func([]interface{}) (interface{}, error), arity int) error {
	stack, err := GetOperandStack(ctx)
	if err != nil {
		return err
	}

	if err := stack.Do(op, arity); err != nil {
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

func loadOp(op func(context.Context, interface{}) (interface{}, error)) Process {
	return func(ctx context.Context, imms []interface{}) error {
		var v interface{}
		var err error
		if len(imms) > 0 {
			v = imms[0]
		} else {
			v, err = Pop(ctx)
		}
		if err != nil {
			return err
		}

		v, err = op(ctx, v)
		if err != nil {
			return err
		}

		if err := Push(ctx, v); err != nil {
			return err
		}

		GetPC(ctx).Increment()
		return nil
	}
}

func ld(ctx context.Context, v interface{}) (interface{}, error) {
	v, _ = GetGlobalHeap(ctx).Load(ToString(v))
	return v, nil
}

func lda(ctx context.Context, v interface{}) (interface{}, error) {
	return GetArgument(ctx, ToInteger(v))
}

func ldl(ctx context.Context, v interface{}) (interface{}, error) {
	lh, err := GetLocalHeap(ctx)
	if err != nil {
		return nil, err
	}

	v, _ = lh.Load(ToString(v))
	return v, nil
}

func storeOp(op func(context.Context, []interface{}) error) Process {
	return func(ctx context.Context, imms []interface{}) error {
		var vs []interface{}
		var err error
		switch len(imms) {
		case 0:
			vs, err = MultiPop(ctx, 2)
		case 1:
			vs, err = MultiPop(ctx, 1)
			if err == nil {
				vs = append(vs, imms[0])
			}
		default:
			vs = imms
		}
		if err != nil {
			return err
		}

		if err := op(ctx, vs); err != nil {
			return err
		}

		GetPC(ctx).Increment()
		return nil
	}
}

func st(ctx context.Context, vs []interface{}) error {
	GetGlobalHeap(ctx).Store(ToString(vs[0]), vs[1])
	return nil
}

func stl(ctx context.Context, vs []interface{}) error {
	lh, err := GetLocalHeap(ctx)
	if err != nil {
		return err
	}

	lh.Store(ToString(vs[0]), vs[1])
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

	frame := newFrame()
	frame.Arguments = argv
	frame.ReturnTo.SetValue(pc.GetValue())
	getCallStack(ctx).Push(frame)

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

	frame, err := getFrame(ctx)
	if err != nil {
		return err
	}

	GetPC(ctx).SetValue(frame.ReturnTo.GetValue())
	_, err = getCallStack(ctx).Pop()
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

	if ToBoolean(v) {
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

	if !ToBoolean(v) {
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
	return BooleanValue(Equal(vs[0], vs[1])), nil
}

func ne(vs []interface{}) (interface{}, error) {
	return BooleanValue(!Equal(vs[0], vs[1])), nil
}

func gt(vs []interface{}) (interface{}, error) {
	return BooleanValue(Less(vs[1], vs[0])), nil
}

func ge(vs []interface{}) (interface{}, error) {
	return BooleanValue(Less(vs[1], vs[0]) || Equal(vs[1], vs[0])), nil
}

func lt(vs []interface{}) (interface{}, error) {
	return BooleanValue(Less(vs[0], vs[1])), nil
}

func le(vs []interface{}) (interface{}, error) {
	return BooleanValue(Less(vs[0], vs[1]) || Equal(vs[0], vs[1])), nil
}

func not(vs []interface{}) (interface{}, error) {
	return BooleanValue(!ToBoolean(vs[0])), nil
}

func and(vs []interface{}) (interface{}, error) {
	return BooleanValue(ToBoolean(vs[0]) && ToBoolean(vs[1])), nil
}

func or(vs []interface{}) (interface{}, error) {
	return BooleanValue(ToBoolean(vs[0]) || ToBoolean(vs[1])), nil
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

func loadStoreOp(op func(ctx context.Context, v interface{}) error) Process {
	return func(ctx context.Context, imms []interface{}) error {
		var v interface{}
		var err error
		if len(imms) > 0 {
			v = imms[0]
		} else {
			v, err = Pop(ctx)
		}
		if err != nil {
			return err
		}

		if err := op(ctx, v); err != nil {
			return err
		}

		GetPC(ctx).Increment()
		return nil
	}
}

func inc(ctx context.Context, v interface{}) error {
	h := GetGlobalHeap(ctx)
	k := ToString(v)
	v, _ = h.Load(k)
	h.Store(k, ToNumber(v)+1.0)
	return nil
}

func incl(ctx context.Context, v interface{}) error {
	lh, err := GetLocalHeap(ctx)
	if err != nil {
		return err
	}

	k := ToString(v)
	v, _ = lh.Load(k)
	lh.Store(k, ToNumber(v)+1.0)
	return nil
}

func dec(ctx context.Context, v interface{}) error {
	h := GetGlobalHeap(ctx)
	k := ToString(v)
	v, _ = h.Load(k)
	h.Store(k, ToNumber(v)-1.0)
	return nil
}

func decl(ctx context.Context, v interface{}) error {
	lh, err := GetLocalHeap(ctx)
	if err != nil {
		return err
	}

	k := ToString(v)
	v, _ = lh.Load(k)
	lh.Store(k, ToNumber(v)-1.0)
	return nil
}
