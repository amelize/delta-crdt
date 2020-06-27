package kernel

import (
	"log"
	"reflect"
)

type ContextData struct {
	CasualContext map[string]int64
	Cloud         []Pair
}

type DotContext struct {
	casualContext map[string]int64
	dotCloud      map[Pair]bool
}

func (ctx DotContext) GetData() ContextData {
	cloud := make([]Pair, 0, len(ctx.dotCloud))

	for pair := range ctx.dotCloud {
		cloud = append(cloud, pair)
	}

	data := ContextData{
		CasualContext: ctx.casualContext,
		Cloud:         cloud,
	}

	return data
}

func (ctx DotContext) Copy() *DotContext {
	cp := NewDotContext()
	for k, v := range ctx.casualContext {
		cp.casualContext[k] = v
	}

	for k, v := range ctx.dotCloud {
		cp.dotCloud[k] = v
	}

	return cp
}

func NewFromData(data ContextData) *DotContext {
	dotCloud := make(map[Pair]bool)

	for _, pair := range data.Cloud {
		dotCloud[pair] = true
	}

	return &DotContext{
		casualContext: data.CasualContext,
		dotCloud:      dotCloud,
	}
}

func NewDotContext() *DotContext {
	return &DotContext{
		casualContext: make(map[string]int64),
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
		log.Printf("rejoin")
		return
	}

	itOne := NewIterator(dotKernel.Dots)
	itTwo := NewIterator(other.Dots)

	for itOne.HasMore() || itTwo.HasMore() {
		//
		if itOne.HasMore() && (!itTwo.HasMore() || pairCompare(itOne.Key().(Pair), itTwo.Key().(Pair))) {
			p := itOne.Key().(Pair)

			log.Printf("1: %+v", p)

			itOne.Next()

			// if other.Ctx.dotin(p) {
			// 	dotKernel.Dots.Remove(p)
			// }
		} else if itTwo.HasMore() && (!itOne.HasMore() || pairCompare(itTwo.Key().(Pair), itOne.Key().(Pair))) {
			p := itTwo.Key().(Pair)

			log.Printf("2: %+v %+v", p, itTwo.Value())

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
