package encoding

import (
	"github.com/delta-crdt/kernel"
)

type ContextData struct {
	CasualContext map[string]int32
	DotCloud      map[kernel.Pair]bool
}

type CRDTData struct {
	Context ContextData
	Data    map[interface{}]interface{}
}

type Encoder interface {
	Encode(kernel *CRDTData) ([]byte, error)
}
