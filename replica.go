package crdt

import (
	"time"

	"github.com/amelize/delta-crdt/broadcaster"
)

// Replica Replica define replica instance
type Replica struct {
	replicaID string
	objects   *broadcaster.Objects
}

// NewReplica creates a new replica with a specific ID. The ID must be unique cluster-wide.
func NewReplica(replicaID string) *Replica {
	replica := &Replica{
		replicaID: replicaID,
		objects:   broadcaster.NewObjects(),
	}

	go replica.loop()

	return replica
}

func (replica *Replica) Update(name string, data interface{}) error {
	updateFunction, _ := replica.objects.Get(name).Update(data)

	return updateFunction()
}

// CreateNewAWORSet creates new AWORSet inside replica with a specific name.
func (replica *Replica) CreateNewAWORSet(name string, handler AworsetBroadcastHandler) *Aworset {
	set := NewAworset(replica.replicaID, handler)
	replica.objects.Add(name, set)

	return set
}

func (repl *Replica) broadcast() {
	// need to replace with Visitor pattern
	head := repl.objects.GetChangedHead()

	for head != "" {
		obj := repl.objects.Get(head)
		fun, _ := obj.Broadcast(repl.replicaID, head)

		err := fun()
		if err != nil {
			repl.objects.Resend(head)
		}

		head = repl.objects.GetChangedHead()
	}
}

func (b *Replica) loop() {
	// TODO: stop
	ticker := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-ticker.C:
			b.broadcast()

		}
	}
}
