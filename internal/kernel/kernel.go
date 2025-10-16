package kernel

import (
	"reflect"
)

type DotKernel struct {
	// I use RB Three instead of map. map[Pair]interface{}
	Dots *RBTree
	Ctx  *DotContext
}

func NewDotKernel() *DotKernel {
	ctx := NewDotContext()
	return &DotKernel{
		Dots: NewTreeMap(lessPair, equalPair), // make(map[Pair]interface{}),
		Ctx:  ctx,
	}
}

func NewDotKernelWithContext(context *DotContext) *DotKernel {
	return &DotKernel{
		Dots: NewTreeMap(lessPair, equalPair), // make(map[Pair]interface{}),
		Ctx:  context,
	}
}

func (dotKernel DotKernel) Add(id int64, value interface{}) *DotKernel {
	// new result
	result := NewDotKernel()

	// create new dot
	dot := dotKernel.Ctx.makeDot(id)
	dotKernel.Dots.Insert(dot, value)

	// make delta
	result.Dots.Insert(dot, value)
	result.Ctx.insertDot(dot, true)

	return result
}

func (dotKernel DotKernel) RemoveValue(value interface{}) *DotKernel {
	res := NewDotKernel()

	iterator := NewIterator(dotKernel.Dots)
	for iterator.HasMore() {
		if reflect.DeepEqual(iterator.Value(), value) {
			k := iterator.Key().(Pair)
			res.Ctx.insertDot(k, false)
			iterator.Next()
			dotKernel.Dots.Remove(k)
		} else {
			iterator.Next()
		}
	}

	res.Ctx.compact()

	return res
}

func (dotKernel DotKernel) RemovePair(value Pair) *DotKernel {
	res := NewDotKernel()

	exists := dotKernel.Dots.Exists(value)
	if exists {
		res.Ctx.insertDot(value, false)
		dotKernel.Dots.Remove(value)
	}

	res.Ctx.compact()

	return res
}

func (dotKernel DotKernel) RemoveAll() *DotKernel {

	res := NewDotKernel()
	iterator := NewIterator(dotKernel.Dots)
	for iterator.HasMore() {
		k := iterator.Key().(Pair)
		res.Ctx.insertDot(k, false)
		iterator.Next()

		dotKernel.Dots.Remove(k)
	}

	res.Ctx.compact()
	return res
}

func (dotKernel *DotKernel) Join(other *DotKernel) {
	if dotKernel == other {
		return
	}

	itOne := dotKernel.Dots.GetIterator()
	itTwo := other.Dots.GetIterator()

	for itOne.HasMore() || itTwo.HasMore() {
		//
		if itOne.HasMore() && (!itTwo.HasMore() || pairCompare(itOne.Key().(Pair), itTwo.Key().(Pair))) {
			p := itOne.Key().(Pair)

			// dot in other replica
			if other.Ctx.dotin(p) {
				dotKernel.Dots.Remove(p)
			}

			itOne.Next()
		} else if itTwo.HasMore() && (!itOne.HasMore() || pairCompare(itTwo.Key().(Pair), itOne.Key().(Pair))) {
			p := itTwo.Key().(Pair)

			if !dotKernel.Ctx.dotin(p) {
				dotKernel.Dots.Insert(p, itTwo.Value())
			}

			itTwo.Next()
		} else if itOne.HasMore() && itTwo.HasMore() {
			itOne.Next()
			itTwo.Next()
		}
	}

	dotKernel.Ctx.Join(other.Ctx)
}
