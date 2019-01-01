package kernel

type Crdt interface {
	Context() *DotContext
	Join(interface{})
}

type Resetable interface {
	Crdt
	Reset() Resetable
}

// Embedable Uses by types that embeds other types
type Embedable interface {
	Resetable
}
