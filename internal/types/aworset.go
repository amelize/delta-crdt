package types

import (
	"errors"
	"log"
	"sync"

	"github.com/amelize/delta-crdt/internal/aworset"
	"github.com/amelize/delta-crdt/internal/broadcaster"
)

var NoBroadcastHandler = errors.New("No broadcast handler")

type AworsetBroadcastHandler interface {
	Broadcast(replicaID int64, name string, aworset *aworset.AWORSet) error
	OnUpdate(data interface{}) (*aworset.AWORSet, error)
}

type Aworset struct {
	set       *aworset.AWORSet
	deltas    *aworset.AWORSet
	lock      sync.RWMutex
	broadcast AworsetBroadcastHandler
	name      string

	onChangeHandler broadcaster.ChangeHandlerInterface
	onUpdated       broadcaster.UpdatedHandlerInterface
}

func NewAworset(replica int64, name string, broadcastHandler AworsetBroadcastHandler) *Aworset {
	return &Aworset{
		set:       aworset.New(replica),
		deltas:    aworset.NewForDelta(),
		broadcast: broadcastHandler,
		name:      name,
	}
}

func (awset *Aworset) GetName() string {
	return awset.name
}

func (awset *Aworset) SetOnChanged(handler broadcaster.ChangeHandlerInterface) {
	awset.onChangeHandler = handler
}

func (awset *Aworset) SetOnUpdated(onUpdated broadcaster.UpdatedHandlerInterface) {
	awset.onUpdated = onUpdated
}

func (awset *Aworset) SetBroadcastHandler(broadcast AworsetBroadcastHandler) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.broadcast = broadcast
}

func (awset *Aworset) Add(val interface{}) {
	awset.lock.Lock()

	// Notify about change
	defer func() {
		if awset.onChangeHandler != nil {
			awset.onChangeHandler.OnChange(awset.name)
		}

		awset.lock.Unlock()
	}()

	awset.deltas.Join(awset.set.Add(val))
	awset.deltas.Dump()
}

func (awset *Aworset) Remove(val interface{}) {
	awset.lock.Lock()
	defer func() {
		if awset.onChangeHandler != nil {
			awset.onChangeHandler.OnChange(awset.name)
		}

		awset.lock.Unlock()
	}()

	change := awset.set.Remove(val)
	awset.deltas.Join(change)
}

func (awset *Aworset) Reset() {
	awset.lock.Lock()
	defer func() {
		if awset.onChangeHandler != nil {
			awset.onChangeHandler.OnChange(awset.name)
		}

		awset.lock.Unlock()
	}()

	change := awset.set.Reset().(interface{})

	if awset.deltas != nil {
		awset.deltas.Join(change)
	} else {
		awset.deltas = change.(*aworset.AWORSet)
	}
}

func (awset *Aworset) Value() map[interface{}]bool {
	awset.lock.RLock()
	defer awset.lock.RUnlock()

	return awset.set.Value()
}

func (awset *Aworset) In(val interface{}) bool {
	awset.lock.RLock()
	defer awset.lock.RUnlock()

	return awset.set.In(val)
}

// GetChanges returns changes for broadcast and clears changes.
func (awset *Aworset) Broadcast(replicaID int64, name string) error {
	awset.lock.Lock()
	defer func() {
		awset.lock.Unlock()
	}()

	result := awset.deltas

	if awset.broadcast == nil {
		return NoBroadcastHandler
	}

	handler := awset.broadcast
	if result == nil {
		return nil
	}

	log.Printf("broadcast values")
	err := handler.Broadcast(replicaID, name, result)
	if err != nil {
		return err
	}

	log.Printf("broadcast done")
	return nil
}

func (awset *Aworset) Update(data interface{}) error {
	awset.lock.RLock()
	defer func() {
		awset.lock.RUnlock()
	}()

	if awset.broadcast == nil {
		return NoBroadcastHandler
	}

	set, err := awset.broadcast.OnUpdate(data)
	if err != nil {
		return err
	}

	awset.set.Join(set)

	if awset.onUpdated != nil {
		go awset.onUpdated.OnUpdate()
	}

	return nil
}

// Join joins broadcasted changes into set
func (awset *Aworset) Join(data interface{}) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.set.Join(data)
}
