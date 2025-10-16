package aworset

import (
	"log"
	"reflect"

	"github.com/amelize/delta-crdt/internal/kernel"
)

type AWORSet struct {
	id        int64
	dotKernel *kernel.DotKernel
}

func New(id int64) *AWORSet {
	return &AWORSet{
		id:        id,
		dotKernel: kernel.NewDotKernel(),
	}
}

func NewForDelta() *AWORSet {
	return &AWORSet{
		id:        0,
		dotKernel: kernel.NewDotKernel(),
	}
}

func NewFromKernel(kernel *kernel.DotKernel) *AWORSet {
	return &AWORSet{
		dotKernel: kernel,
	}
}

func NewWithContext(id int64, ctx *kernel.DotContext) *AWORSet {
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

// GetKernel retutns kernel
func (set AWORSet) GetKernel() *kernel.DotKernel {
	return set.dotKernel
}

// Context returns content
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

	delta := set.dotKernel.Add(set.id, val)

	res.dotKernel.Join(delta)

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

func (set *AWORSet) Join(other interface{}) {
	set.Dump()

	otherDot, ok := other.(*AWORSet)
	if ok {
		set.dotKernel.Join(otherDot.dotKernel)

	} else {
		panic("wrong type")
	}
}

func (currentKernel *AWORSet) Dump() {
	it := currentKernel.dotKernel.Dots.GetIterator()

	log.Printf("Dump start ----- ")
	for it.HasMore() {
		log.Printf("\t%+v %+v", it.Key(), it.Value())
		it.Next()
	}

	log.Printf("Dump end ----- ")
}
