package ormap

import (
	"github.com/delta-crdt/aworset"
	"github.com/delta-crdt/ccounter"
	"github.com/delta-crdt/kernel"
)

type NewItem = func(id string, ctx *kernel.DotContext) kernel.Embedable

type ORMap struct {
	id      string
	data    *kernel.RBTree
	ctx     *kernel.DotContext
	less    kernel.Less
	equal   kernel.Equal
	newFunc NewItem
}

func New(id string, less kernel.Less, equal kernel.Equal, newFunc NewItem) *ORMap {
	return &ORMap{
		id:      id,
		data:    kernel.New(less, equal),
		ctx:     kernel.NewDotContext(),
		less:    less,
		equal:   equal,
		newFunc: newFunc,
	}
}

func new(less kernel.Less, equal kernel.Equal, newFunc NewItem) *ORMap {
	return &ORMap{
		data:    kernel.New(less, equal),
		ctx:     kernel.NewDotContext(),
		less:    less,
		equal:   equal,
		newFunc: newFunc,
	}
}

func (ormap ORMap) Context() *kernel.DotContext {
	return ormap.ctx
}

func NewWithAworset(id string, less kernel.Less, equal kernel.Equal) *ORMap {
	return New(id, less, equal, func(id string, ctx *kernel.DotContext) kernel.Embedable {
		return aworset.NewWithContext(id, ctx)
	})
}

func NewWithStingKey(id string, newFunc NewItem) *ORMap {
	return New(id, kernel.StringLess, kernel.StringEqual, newFunc)
}

func NewWithAworsetStringKey(id string) *ORMap {
	return NewWithAworset(id, kernel.StringLess, kernel.StringEqual)
}

func (ormap *ORMap) Get(key interface{}) interface{} {
	value := ormap.data.Get(key)
	if value == nil {
		empty := ormap.newFunc(ormap.id, ormap.ctx)

		ormap.data.Insert(key, empty)
		return empty
	}

	return value
}

func (ormap *ORMap) Erase(key interface{}) kernel.Resetable {
	res := new(ormap.less, ormap.equal, ormap.newFunc)

	val := ormap.data.Get(key)
	if val != nil {
		v := val.(kernel.Embedable).Reset()
		res.ctx.Join(v.Context())
		ormap.data.Remove(key)
	}

	return res
}

func (ormap *ORMap) Reset() kernel.Resetable {
	res := new(ormap.less, ormap.equal, ormap.newFunc)

	if !ormap.data.Empty() {
		mit := kernel.NewIterator(ormap.data)

		for mit.HasMore() {
			v := mit.Value().(kernel.Embedable).Reset()
			res.ctx.Join(v.Context())
			mit.Next()
		}

		ormap.data.Clear()
	}

	return res
}

func (ormap *ORMap) Join(other interface{}) {
	otherMap, ok := other.(*ORMap)
	if ok {
		ormap.join(otherMap)
	} else {
		panic("wrong type")
	}
}

func (ormap *ORMap) join(other *ORMap) {
	imctx := ormap.ctx.Copy()

	mit := kernel.NewIterator(ormap.data)
	mito := kernel.NewIterator(other.data)

	for mit.HasMore() && mito.HasMore() {
		if mit.HasMore() && (!mito.HasMore() || ormap.less(mit.Key(), mito.Key())) {
			empty := ormap.newFunc(ormap.id, other.ctx)
			mit.Value().(kernel.Embedable).Join(empty)

			ormap.ctx = imctx

			mit.Next()
		} else if mito.HasMore() && (!mit.HasMore() || ormap.less(mito.Key(), mit.Key())) {
			val := ormap.data.Get(mito.Key()).(kernel.Embedable)
			val.Join(mito.Value())

			ormap.ctx = imctx

			mito.Next()
		} else if mito.HasMore() && mit.HasMore() {
			val := ormap.data.Get(mit.Key()).(kernel.Embedable)
			val.Join(mito.Value())

			ormap.ctx = imctx

			mito.Next()
			mit.Next()
		}
	}

	ormap.ctx.Join(other.ctx)
}

func (ormap *ORMap) GetAsAworSet(id string) *aworset.AWORSet {
	return ormap.Get(id).(*aworset.AWORSet)
}

func (ormap *ORMap) GetAsIntCounter(id string) *ccounter.IntCounter {
	return ormap.Get(id).(*ccounter.IntCounter)
}

func IntCounter(id string, ctx *kernel.DotContext) kernel.Embedable {
	return ccounter.NewIntCounterWithContex(id, ctx)
}
