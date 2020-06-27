package crdt

import (
	"time"

	"github.com/amelize/delta-crdt/broadcaster"
	"github.com/google/uuid"
)

// Replica Replica define replica instance
type Replica struct {
	replicaID     string
	objects       *broadcaster.Objects
	broadcastRate time.Duration
}

func NewReplicaWithSelfUniqueAddress(broadcastRate time.Duration) *Replica {
	return NewReplica(uuid.New().String(), broadcastRate)
}

// NewReplica creates a new replica with a specific ID. The ID must be unique cluster-wide.
func NewReplica(replicaID string, broadcastRate time.Duration) *Replica {
	replica := &Replica{
		replicaID:     replicaID,
		objects:       broadcaster.NewObjects(),
		broadcastRate: broadcastRate,
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

func (replica *Replica) CreateCCounter(name string, handler CCounterBroadcastHandler) *CCounter {
	counter := NewCCounter(replica.replicaID, handler)
	replica.objects.Add(name, counter)

	return counter
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

func (b Replica) loop() {
	// TODO: stop
	ticker := time.NewTicker(b.broadcastRate)

	for {
		select {
		case <-ticker.C:
			b.broadcast()
		}
	}
}
