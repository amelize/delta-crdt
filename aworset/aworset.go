package aworset

import (
	"reflect"

	"github.com/delta-crdt/kernel"
)

type AWORSet struct {
	id        string
	dotKernel *kernel.DotKernel
}

func New(id string) *AWORSet {
	return &AWORSet{
		id:        id,
		dotKernel: kernel.NewDotKernel(),
	}
}

func NewFromKernel(kernel *kernel.DotKernel) *AWORSet {
	return &AWORSet{
		dotKernel: kernel,
	}
}

func NewWithContext(id string, ctx *kernel.DotContext) *AWORSet {
	return &AWORSet{
		id:        id,
		dotKernel: kernel.NewDotKernelWithContext(ctx),
	}
}

func new() *AWORSet {
	return &AWORSet{
		dotKernel: kernel.NewDotKernel(),
	}
}

func (set AWORSet) GetKernel() *kernel.DotKernel {
	return set.dotKernel
}

func (set AWORSet) Context() *kernel.DotContext {
	return set.dotKernel.Ctx
}

func (set AWORSet) Value() map[interface{}]bool {
	result := make(map[interface{}]bool)

	it := kernel.NewIterator(set.dotKernel.Dots)
	for it.HasMore() {
		result[it.Value()] = true
		it.Next()
	}

	return result
}

func (set AWORSet) In(val interface{}) bool {
	it := kernel.NewIterator(set.dotKernel.Dots)
	for it.HasMore() {
		if reflect.DeepEqual(it.Value(), val) {
			return true
		}

		it.Next()
	}
	return false
}

func (set AWORSet) Add(val interface{}) *AWORSet {
	res := new()

	res.dotKernel = set.dotKernel.RemoveValue(val)
	res.dotKernel.Join(set.dotKernel.Add(set.id, val))

	return res
}

func (set AWORSet) Remove(val interface{}) *AWORSet {
	res := new()

	res.dotKernel = set.dotKernel.RemoveValue(val)
	return res
}

func (set AWORSet) Reset() kernel.Resetable {
	res := new()

	res.dotKernel = set.dotKernel.RemoveAll()
	return res
}

func (set AWORSet) Join(other interface{}) {
	otherDot, ok := other.(*AWORSet)
	if ok {
		set.dotKernel.Join(otherDot.dotKernel)
	} else {
		panic("wrong type")
	}
}
