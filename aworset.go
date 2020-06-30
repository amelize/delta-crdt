package crdt

import (
	"errors"
	"log"
	"sync"

	"github.com/amelize/delta-crdt/aworset"
	"github.com/amelize/delta-crdt/broadcaster"
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

	onChanged broadcaster.OnChanged
	onUpdated broadcaster.OnUpdated
}

func NewAworset(replica int64, broadcastHandler AworsetBroadcastHandler) *Aworset {
	return &Aworset{
		set:       aworset.New(replica),
		deltas:    aworset.NewForDelta(),
		broadcast: broadcastHandler,
	}
}

func (awset *Aworset) SetOnChanged(onChanged broadcaster.OnChanged) {
	awset.onChanged = onChanged
}

func (awset *Aworset) SetOnUpdated(onUpdated broadcaster.OnUpdated) {
	awset.onUpdated = onUpdated
}

func (awset *Aworset) SetBroadcastHandler(broadcast AworsetBroadcastHandler) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.broadcast = broadcast
}

func (awset *Aworset) Add(val interface{}) {
	awset.lock.Lock()
	defer func() {
		awset.lock.Unlock()

		if awset.onChanged != nil {
			go awset.onChanged()
		}
	}()

	awset.deltas.Join(awset.set.Add(val))
	awset.deltas.Dump()

	log.Printf("---------------------------------------- END")
}

func (awset *Aworset) Remove(val interface{}) {
	awset.lock.Lock()
	defer func() {
		awset.lock.Unlock()

		if awset.onChanged != nil {
			go awset.onChanged()
		}
	}()

	change := awset.set.Remove(val)

	awset.deltas.Join(change)
}

func (awset *Aworset) Reset() {
	awset.lock.Lock()
	defer func() {
		awset.lock.Unlock()

		if awset.onChanged != nil {
			go awset.onChanged()
		}
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
func (awset *Aworset) Broadcast(replicaID int64, name string) (broadcaster.SendFunction, error) {
	awset.lock.Lock()
	defer func() {
		awset.lock.Unlock()
	}()

	result := awset.deltas

	if awset.broadcast == nil {
		return nil, NoBroadcastHandler
	}

	handler := awset.broadcast
	// awset.result = nil

	sendFunction := func() error {
		if result == nil {
			return nil
		}

		log.Printf("broadcast values")
		err := handler.Broadcast(replicaID, name, result)
		if err != nil {
			awset.lock.Lock()
			defer awset.lock.Unlock()

			// if awset.deltas != nil {
			// 	result.Join(awset.deltas)
			// }
			// awset.deltas = result

			return err
		}

		return nil
	}

	return sendFunction, nil
}

func (awset *Aworset) Update(data interface{}) (broadcaster.UpdateFunction, error) {
	awset.lock.RLock()
	defer func() {
		awset.lock.RUnlock()
	}()

	if awset.broadcast == nil {
		return nil, NoBroadcastHandler
	}

	updateFunction := func() error {
		set, err := awset.broadcast.OnUpdate(data)
		if err != nil {
			return err
		}

		awset.Join(set)

		if awset.onUpdated != nil {
			awset.onUpdated()
		}

		return nil
	}

	return updateFunction, nil
}

// Join joins broadcasted changes into set
func (awset *Aworset) Join(data interface{}) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.set.Join(data)
}
