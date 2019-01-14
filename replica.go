package crdt

import (
	"github.com/delta-crdt/encoding"
	"github.com/delta-crdt/kernel"
)

// Sender interface that you need to implement to make replica send broadcasts.
type Sender interface {
	Sender(name string, data []byte) error
}

// Replica Replica define replica instance
type Replica struct {
	replicaID string
	data      map[string]kernel.Joinable
}

// NewReplica creates a new replica with a specific ID. The ID must be unique cluster-wide.
func NewReplica(replicaID string) *Replica {
	return &Replica{
		replicaID: replicaID,
		data:      make(map[string]kernel.Joinable),
	}
}

func (replica *Replica) Update(name string, data []byte) error {

}

// CreateNewAWORSet creates new AWORSet inside replica with a specific name.
func (replica *Replica) CreateNewAWORSet(name string, encoder encoding.Encoder) *Aworset {
	set := NewAworset(replica.replicaID)
	replica.data[name] = set

	return set
}
