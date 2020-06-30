package kernel

import "sort"

type Pair struct {
	First  int64
	Second uint64
}

func NewPair(first int64, second uint64) Pair {
	return Pair{First: first, Second: second}
}

func (this *Pair) Compare(other Pair) bool {
	if this.First < other.First {
		return true
	}

	if this.First > other.First {
		return false
	}

	if this.Second < other.Second {
		return true
	}

	return false
}

func (this *Pair) Equal(other Pair) bool {
	return this.First == other.First && this.Second == other.Second
}

func pairCompare(a Pair, b Pair) bool {
	return a.Compare(b)
}

func firstLess(a Pair, b Pair) bool {
	return a.First < b.First
}

func pairCompareEqual(a Pair, b Pair) bool {
	return a.Equal(b)
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

	return ai.Compare(aj)
}

func lessPair(a interface{}, b interface{}) bool {
	return pairCompare(a.(Pair), b.(Pair))
}

func equalPair(a interface{}, b interface{}) bool {
	return pairCompareEqual(a.(Pair), b.(Pair))
}

func StringLess(a interface{}, b interface{}) bool {
	return a.(string) < b.(string)
}

func StringEqual(a interface{}, b interface{}) bool {
	return a.(string) == b.(string)
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

	return ai.Compare(aj)
}

type casualContextIterator struct {
	values  orderedPair
	current int
}

func CreateCCIterator(source map[int64]uint64) casualContextIterator {
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
