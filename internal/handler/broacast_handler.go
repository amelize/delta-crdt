package handler

type BroadcastHandlerInterface interface {
	Broadcast(replicaID int64, name string, crdt interface{})
}

// Converts data to crdt
type SerializeInterface interface {
	Deserialize(data interface{}) (interface{}, error)
}
