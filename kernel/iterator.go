package kernel

import "sort"

type Pair struct {
	First  string
	Second int32
}

func pairCompair(a Pair, b Pair) bool {
	if a.First < b.First {
		return true
	}

	if a.First > b.First {
		return false
	}

	if a.Second < b.Second {
		return true
	}

	return false
}

type dot struct {
	pair  Pair
	value interface{}
}

type orderedDots []dot

func (a orderedDots) Len() int      { return len(a) }
func (a orderedDots) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a orderedDots) Less(i, j int) bool {
	ai := a[i].pair
	aj := a[j].pair

	return pairCompair(ai, aj)
}

func createOrderedDots(source map[Pair]interface{}) orderedDots {
	dots := make(orderedDots, 0, len(source))
	for k, v := range source {
		dots = append(dots, dot{pair: k, value: v})
	}

	sort.Sort(dots)

	return dots
}

type dotsIterator struct {
	dots    orderedDots
	current int
}

func newOrderedIterator(source map[Pair]interface{}) dotsIterator {
	return dotsIterator{dots: createOrderedDots(source), current: 0}
}

func (it dotsIterator) hasMore() bool {
	return it.current < len(it.dots)
}

func (it *dotsIterator) next() {
	it.current++
}

func (it dotsIterator) val() dot {
	return it.dots[it.current]
}

type orderedPair []Pair

func (a orderedPair) Len() int      { return len(a) }
func (a orderedPair) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a orderedPair) Less(i, j int) bool {
	ai := a[i]
	aj := a[j]

	return pairCompair(ai, aj)
}

type casualContextIterator struct {
	values  orderedPair
	current int
}

func CreateCCIterator(source map[string]int32) casualContextIterator {
	vals := make(orderedPair, 0, len(source))

	for k, v := range source {
		vals = append(vals, Pair{First: k, Second: v})
	}

	sort.Sort(vals)

	return casualContextIterator{values: vals}
}

func (it casualContextIterator) hasMore() bool {
	return it.current < len(it.values)
}

func (it *casualContextIterator) next() {
	it.current++
}

func (it casualContextIterator) val() Pair {
	return it.values[it.current]
}
