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

func new() *AWORSet {
	return &AWORSet{
		dotKernel: kernel.NewDotKernel(),
	}
}

func (set AWORSet) Value() map[interface{}]bool {
	result := make(map[interface{}]bool)

	for _, v := range set.dotKernel.Dots {
		result[v] = true
	}

	return result
}

func (set AWORSet) In(val interface{}) bool {
	for _, v := range set.dotKernel.Dots {
		if reflect.DeepEqual(val, v) {
			return true
		}
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

func (set AWORSet) Join(other *AWORSet) {
	set.dotKernel.Join(other.dotKernel)
}
