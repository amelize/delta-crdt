package crdt

import (
	"testing"

	"github.com/delta-crdt/aworset"
)

type DummyHandler struct {
	other *Replica
}

func (handler DummyHandler) Broadcast(replicaID, name string, aworset *aworset.AWORSet) error {
	return handler.other.Update(name, aworset)
}

func (handler DummyHandler) OnUpdate(data interface{}) (*aworset.AWORSet, error) {
	return data.(*aworset.AWORSet), nil
}

func TestReplica_CreateNewAWORSet(t *testing.T) {
	lock := make(chan struct{})

	replicaOne := NewReplica("a")
	replicaTwo := NewReplica("b")

	firstHandler := DummyHandler{other: replicaTwo}
	secondHandler := DummyHandler{other: replicaOne}

	setOne := replicaOne.CreateNewAWORSet("user.set", firstHandler)
	setTwo := replicaTwo.CreateNewAWORSet("user.set", secondHandler)

	setTwo.SetOnUpdated(func() {
		lock <- struct{}{}
	})

	setOne.Add("HelloBadge")
	setTwo.Add("WelocomeBadge")

	<-lock
}
