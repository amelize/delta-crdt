package ccounter

import (
	"github.com/amelize/delta-crdt/kernel"
)

type CounterValue interface {
	Base() CounterValue
	Sum(a CounterValue, b CounterValue) CounterValue
	Sub(a CounterValue, b CounterValue) CounterValue
	Max(a CounterValue, b CounterValue) CounterValue
	Value() interface{}
}

type CCounter struct {
	id        string
	dotKernel *kernel.DotKernel
}

func New(id string) *CCounter {
	return &CCounter{
		id:        id,
		dotKernel: kernel.NewDotKernel(),
	}
}

func NewWithContext(id string, ctx *kernel.DotContext) *CCounter {
	return &CCounter{
		id:        id,
		dotKernel: kernel.NewDotKernelWithContext(ctx),
	}
}

func NewWithKernel(id string, kernelData *kernel.DotKernel) *CCounter {
	return &CCounter{
		id:        id,
		dotKernel: kernelData,
	}
}

func new() *CCounter {
	return &CCounter{
		dotKernel: kernel.NewDotKernel(),
	}
}

func (counter CCounter) GetKernel() *kernel.DotKernel {
	return counter.dotKernel
}

func (counter CCounter) Context() *kernel.DotContext {
	return counter.dotKernel.Ctx
}

func (counter CCounter) Inc(val CounterValue) *CCounter {
	res := new()

	dots := make(map[kernel.Pair]bool)
	base := val.Base()

	it := kernel.NewIterator(counter.dotKernel.Dots)

	for it.HasMore() {
		k := it.Key().(kernel.Pair)
		if k.First == counter.id {
			base = val.Max(base, it.Value().(CounterValue))
			dots[k] = true
		}

		it.Next()
	}

	for k := range dots {
		res.dotKernel.Join(counter.dotKernel.RemovePair(k))
	}

	res.dotKernel.Join(counter.dotKernel.Add(counter.id, base.Sum(base, val)))

	return res
}

func (counter CCounter) Dec(val CounterValue) *CCounter {
	res := new()

	dots := make(map[kernel.Pair]bool)
	base := val.Base()

	it := kernel.NewIterator(counter.dotKernel.Dots)

	for it.HasMore() {
		k := it.Key().(kernel.Pair)
		if k.First == counter.id {
			base = val.Max(base, it.Value().(CounterValue))
			dots[k] = true
		}

		it.Next()
	}

	for k := range dots {
		res.dotKernel.Join(counter.dotKernel.RemovePair(k))
	}

	res.dotKernel.Join(counter.dotKernel.Add(counter.id, base.Sub(base, val)))

	return res
}

func (counter CCounter) Value(base CounterValue) CounterValue {
	it := kernel.NewIterator(counter.dotKernel.Dots)
	for it.HasMore() {
		base = base.Sum(base, it.Value().(CounterValue))
		it.Next()
	}

	return base
}

func (counter *CCounter) Reset() kernel.Resetable {
	res := new()

	res.dotKernel = counter.dotKernel.RemoveAll()

	return res
}

func (counter CCounter) Join(other interface{}) {
	otherCounter, ok := other.(*CCounter)
	if ok {
		counter.dotKernel.Join(otherCounter.dotKernel)
	} else {
		panic("wrong type")
	}
}

type IntValue int64

func (c IntValue) Base() CounterValue {
	return IntValue(0)
}

func (c IntValue) Sum(a CounterValue, b CounterValue) CounterValue {
	return IntValue(a.Value().(IntValue) + b.Value().(IntValue))
}

func (c IntValue) Sub(a CounterValue, b CounterValue) CounterValue {
	return IntValue(a.Value().(IntValue) - b.Value().(IntValue))
}
func (c IntValue) Max(a CounterValue, b CounterValue) CounterValue {
	if a.Value().(IntValue) > b.Value().(IntValue) {
		return a
	}
	return b
}
func (c IntValue) Value() interface{} {
	return c
}

type IntCounter struct {
	counter *CCounter
}

func NewIntCounter(id string) *IntCounter {
	return &IntCounter{
		counter: New(id),
	}
}

func NewIntCounterWithContex(id string, ctx *kernel.DotContext) *IntCounter {
	return &IntCounter{
		counter: NewWithContext(id, ctx),
	}
}

func NewIntCounterWithKernel(id string, ctx *kernel.DotKernel) *IntCounter {
	return &IntCounter{
		counter: NewWithKernel(id, ctx),
	}
}

func (c IntCounter) Context() *kernel.DotContext {
	return c.counter.Context()
}

func (c IntCounter) Inc(val int64) *IntCounter {
	cv := IntValue(val)
	rc := c.counter.Inc(cv)
	return &IntCounter{counter: rc}
}

func (c IntCounter) Dec(val int64) *IntCounter {
	cv := IntValue(val)
	rc := c.counter.Dec(cv)
	return &IntCounter{counter: rc}
}

func (c IntCounter) Reset() kernel.Resetable {
	ctxi := c.counter.Reset().(interface{})

	return &IntCounter{counter: ctxi.(*CCounter)}
}

func (c *IntCounter) Join(other interface{}) {
	intC, ok := other.(*IntCounter)
	if ok {
		c.counter.Join(intC.counter)
	}
}

func (c IntCounter) GetCounter() *CCounter {
	return c.counter
}

func (c IntCounter) GetId() string {
	return c.counter.id
}

func (c IntCounter) Value() int64 {
	val := c.counter.Value(IntValue(0))
	return int64(val.Value().(IntValue))
}
