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

// Embeddable Uses by types that embeds other types
type Embeddable interface {
	Resetable
}
