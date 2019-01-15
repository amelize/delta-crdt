package kernel

type Joinable interface {
	Join(interface{})
}

type Crdt interface {
	Joinable
	Context() *DotContext
}

type Resetable interface {
	Crdt
	Reset() Resetable
}

// Embedable Uses by types that embeds other types
type Embedable interface {
	Resetable
}