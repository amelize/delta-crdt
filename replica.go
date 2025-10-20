package crdt

import (
	"context"
	"math/rand"
	"time"

	"log"

	"github.com/amelize/delta-crdt/internal/broadcaster"
	"github.com/amelize/delta-crdt/internal/types"
)

// Replica Replica define replica instances
type Replica struct {
	replicaID     int64
	objects       *broadcaster.Objects
	broadcastRate time.Duration
	ctx           context.Context
}

func NewReplicaWithSelfUniqueAddress(ctx context.Context, broadcastRate time.Duration) *Replica {
	return NewReplica(ctx, rand.Int63(), broadcastRate)
}

// NewReplica creates a new replica with a specific ID. The ID must be unique cluster-wide.
func NewReplica(ctx context.Context, replicaID int64, broadcastRate time.Duration) *Replica {
	replica := &Replica{
		replicaID:     replicaID,
		objects:       broadcaster.NewObjects(),
		broadcastRate: broadcastRate,
		ctx:           ctx,
	}

	go replica.loop()

	return replica
}

func (replica *Replica) Update(name string, data interface{}) error {
	return replica.objects.Get(name).Update(data)
}

// CreateNewAWORSet creates new AWORSet inside replica with a specific name.
func (replica *Replica) CreateNewAWORSet(name string, handler types.AworsetBroadcastHandler) *types.Aworset {
	set := types.NewAworset(replica.replicaID, name, handler)
	replica.objects.Add(name, set)

	return set
}

func (replica *Replica) CreateCCounter(name string, handler types.CCounterBroadcastHandler) *types.CCounter {
	counter := types.NewCCounter(replica.replicaID, name, handler)
	replica.objects.Add(name, counter)

	return counter
}

func (replica *Replica) GetAsCCounter(name string) *types.CCounter {
	return replica.objects.Get(name).(*types.CCounter)
}

func (replica *Replica) GetAsAworSet(name string) *types.Aworset {
	return replica.objects.Get(name).(*types.Aworset)
}

func (replica *Replica) broadcast() {
	head := replica.objects.GetChangedHead()

	for head != nil {
		obj := replica.objects.Get(*head)
		err := obj.Broadcast(replica.replicaID, *head)

		if err != nil {
			replica.objects.Resend(*head)
		}

		head = replica.objects.GetChangedHead()
	}
}

func (b Replica) loop() {
	ticker := time.NewTicker(b.broadcastRate)

	for {
		select {
		case <-ticker.C:
			b.broadcast()
		case <-b.ctx.Done():
			log.Printf("Shutting down")
			return
		}
	}
}
