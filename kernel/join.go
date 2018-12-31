package kernel

// Joinable Uses by types that embeds other types
type Embedable interface {
	Join(interface{})
	Context() *DotContext
}
