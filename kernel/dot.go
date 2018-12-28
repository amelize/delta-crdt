package kernel

import "reflect"

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
	Dots map[Pair]interface{}
	Ctx  *DotContext
}

func NewDotKernel() *DotKernel {
	ctx := NewDotContext()
	return &DotKernel{
		Dots: make(map[Pair]interface{}),
		Ctx:  ctx,
	}
}

func (dotKernel DotKernel) Add(id string, value interface{}) *DotKernel {
	dot := dotKernel.Ctx.makeDot(id)
	dotKernel.Dots[dot] = value

	res := NewDotKernel()
	res.Dots[dot] = value
	res.Ctx.insertDot(dot, true)

	return res
}

func (dotKernel DotKernel) RemoveValue(value interface{}) *DotKernel {
	res := NewDotKernel()

	for k, v := range dotKernel.Dots {
		if reflect.DeepEqual(v, value) {
			res.Ctx.insertDot(k, false)
			delete(dotKernel.Dots, k)
		}
	}

	res.Ctx.compact()

	return res
}

func (dotKernel DotKernel) RemovePair(value Pair) *DotKernel {
	res := NewDotKernel()

	_, exists := dotKernel.Dots[value]
	if exists {
		res.Ctx.insertDot(value, false)
		delete(dotKernel.Dots, value)
	}

	res.Ctx.compact()

	return res
}

func (dotKernel DotKernel) RemoveAll() *DotKernel {
	res := NewDotKernel()
	for k := range dotKernel.Dots {
		res.Ctx.insertDot(k, false)
		delete(dotKernel.Dots, k)
	}

	res.Ctx.compact()
	return res
}

func (dotKernel *DotKernel) Join(other *DotKernel) {
	if dotKernel == other {
		return
	}

	it := newOrderedIterator(dotKernel.Dots)
	ito := newOrderedIterator(other.Dots)

	for it.hasMore() || ito.hasMore() {
		if it.hasMore() && (!ito.hasMore() || pairCompair(it.val().pair, ito.val().pair)) {
			p := it.val().pair
			if other.Ctx.dotin(p) {
				delete(dotKernel.Dots, p)
				it.next()
			} else {
				it.next()
			}
		} else if ito.hasMore() && (!it.hasMore() || pairCompair(ito.val().pair, it.val().pair)) {
			p := ito.val().pair
			if !dotKernel.Ctx.dotin(p) {
				dotKernel.Dots[p] = ito.val().value
			}

			ito.next()
		} else if it.hasMore() && ito.hasMore() {
			it.next()
			ito.next()
		}
	}

	dotKernel.Ctx.Join(other.Ctx)
}
