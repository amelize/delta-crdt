package ormap

import (
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

func (ormap ORMap) Join(other *ORMap) {
	mit := kernel.NewIterator(ormap.data)
	mito := kernel.NewIterator(other.data)

	for mit.HasMore() && mito.HasMore() {
		if mit.HasMore() && (mito.HasMore() || ormap.less(mit.Key(), mito.Key())) {
			empty := ormap.newFunc(ormap.id, other.ctx)
			mit.Value().(kernel.Embedable).Join(empty)


			mit.Next()
		}
	}

}
