package kernel

import (
	"reflect"
)

type DotContext struct {
	casualContext map[string]int32
	dotCloud      map[Pair]bool
}

func NewDotContext() *DotContext {
	return &DotContext{
		casualContext: make(map[string]int32),
		dotCloud:      make(map[Pair]bool),
	}
}

func (ctx DotContext) dotin(p Pair) bool {
	val, ok := ctx.casualContext[p.First]
	if ok {
		if p.Second <= val {
			return true
		}
	}

	if len(ctx.dotCloud) != 0 {
		return true
	}

	return false
}

func (ctx DotContext) compact() {
	needMore := true
	for needMore {
		needMore = false

		for val := range ctx.dotCloud {
			cv, exist := ctx.casualContext[val.First]
			if !exist {
				if val.Second == 1 {
					ctx.casualContext[val.First] = val.Second
					delete(ctx.dotCloud, val)
					needMore = true
				}
			} else {
				if val.Second == cv+1 {
					ctx.casualContext[val.First] = cv + 1
					delete(ctx.dotCloud, val)
					needMore = true
				} else {
					if val.Second <= cv {
						delete(ctx.dotCloud, val)
					}
				}
			}
		}
	}
}

func (ctx DotContext) makeDot(id string) Pair {
	pair := Pair{First: id, Second: 1}
	v, ok := ctx.casualContext[id]
	if ok {
		pair.Second = v + 1
		ctx.casualContext[id] = v + 1
	} else {
		ctx.casualContext[id] = pair.Second
	}

	return pair
}

func (ctx DotContext) insertDot(p Pair, needCompact bool) {
	ctx.dotCloud[p] = true
	if needCompact {
		ctx.compact()
	}
}

func (ctx DotContext) Join(other *DotContext) {
	if &ctx == other {
		return
	}
	it := CreateCCIterator(ctx.casualContext)
	ito := CreateCCIterator(other.casualContext)

	for it.hasMore() || ito.hasMore() {
		if it.hasMore() && (!ito.hasMore() || it.val().First < ito.val().First) {
			it.next()
		} else if ito.hasMore() && !it.hasMore() || ito.val().First < it.val().First {
			pair := ito.val()
			ctx.casualContext[pair.First] = pair.Second
			ito.next()
		} else if it.hasMore() && it.hasMore() {
			cpair := it.val()
			opair := ito.val()
			mx := cpair.Second
			if mx < opair.Second {
				mx = opair.Second
			}

			ctx.casualContext[cpair.First] = mx
			it.next()
			ito.next()
		}
	}

	for k := range other.dotCloud {
		ctx.insertDot(k, false)
	}

	ctx.compact()
}

type DotKernel struct {
	Dots *RBTree //map[Pair]interface{}
	Ctx  *DotContext
}

func NewDotKernel() *DotKernel {
	ctx := NewDotContext()
	return &DotKernel{
		Dots: New(lessPair, equalPair), // make(map[Pair]interface{}),
		Ctx:  ctx,
	}
}

func NewDotKernelWithContext(context *DotContext) *DotKernel {
	return &DotKernel{
		Dots: New(lessPair, equalPair), // make(map[Pair]interface{}),
		Ctx:  context,
	}
}

func (dotKernel DotKernel) Add(id string, value interface{}) *DotKernel {
	dot := dotKernel.Ctx.makeDot(id)
	dotKernel.Dots.Insert(dot, value)

	res := NewDotKernel()
	res.Dots.Insert(dot, value)
	res.Ctx.insertDot(dot, true)

	return res
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

	it := NewIterator(dotKernel.Dots)
	ito := NewIterator(other.Dots)

	for it.HasMore() || ito.HasMore() {
		if it.HasMore() && (!ito.HasMore() || pairCompair(it.Key().(Pair), ito.Key().(Pair))) {
			p := it.Key().(Pair)

			it.Next()

			if other.Ctx.dotin(p) {
				dotKernel.Dots.Remove(p)
			}
		} else if ito.HasMore() && (!it.HasMore() || pairCompair(ito.Key().(Pair), it.Key().(Pair))) {
			p := ito.Key().(Pair)
			if !dotKernel.Ctx.dotin(p) {
				dotKernel.Dots.Insert(p, ito.Value())
			}

			ito.Next()
		} else if it.HasMore() && ito.HasMore() {
			it.Next()
			ito.Next()
		}
	}

	dotKernel.Ctx.Join(other.Ctx)
}
