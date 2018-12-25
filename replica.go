package crdt

// Replica Replica define replica instance
type Replica struct {
	replicaID string
}

func NewReplica(replicaID string) *Replica {
	return &Replica{replicaID: replicaID}
}
