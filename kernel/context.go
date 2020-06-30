package kernel

type ContextData struct {
	CasualContext map[int64]uint64
	Cloud         []Pair
}

type DotContext struct {
	/**
	CausalContext = P(I × N)
		maxi(c) = max({n | (i, n) ∈ c} ∪ {0})
		nexti(c) = (i, maxi(c) + 1)
	*/
	casualContext map[int64]uint64
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
		casualContext: make(map[int64]uint64),
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

	if ctx.countInDotCloud(p) != 0 {
		return true
	}

	return false
}

func (ctx DotContext) countInDotCloud(p Pair) int {
	elementCount := 0

	for k := range ctx.dotCloud {
		if k.Equal(p) {
			elementCount = elementCount + 1
		}
	}

	return elementCount
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

func (ctx DotContext) makeDot(id int64) Pair {
	// Crate new pair
	pair := NewPair(id, 1)

	// is the a exists?
	v, ok := ctx.casualContext[id]
	if ok { // yep, update value
		pair.Second = v + 1
	}

	// store a new dot or update
	ctx.casualContext[id] = pair.Second

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
		} else if ito.hasMore() && (!it.hasMore() || ito.val().First < it.val().First) {
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
