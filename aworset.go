package crdt

import (
	"errors"
	"sync"

	"github.com/delta-crdt/aworset"
	"github.com/delta-crdt/broadcaster"
)

var NoBroadcastHandler = errors.New("No broadcast handler")

type AworsetBroadcastHandler interface {
	Broadcast(replicaID, name string, aworset *aworset.AWORSet) error
	OnUpdate(data interface{}) (*aworset.AWORSet, error)
}

type Aworset struct {
	set       *aworset.AWORSet
	result    *aworset.AWORSet
	lock      sync.RWMutex
	broadcast AworsetBroadcastHandler

	onChanged broadcaster.OnChanged
	onUpdated broadcaster.OnUpdated
}

func NewAworset(replica string, broadcastHandler AworsetBroadcastHandler) *Aworset {
	return &Aworset{
		set:       aworset.New(replica),
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

	change := awset.set.Add(val)

	if awset.result != nil {
		awset.result.Join(change)
	} else {
		awset.result = change
	}
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

	if awset.result != nil {
		awset.result.Join(change)
	} else {
		awset.result = change
	}
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

	if awset.result != nil {
		awset.result.Join(change)
	} else {
		awset.result = change.(*aworset.AWORSet)
	}
}

func (awset *Aworset) Value() interface{} {
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
func (awset *Aworset) Broadcast(replicaID, name string) (broadcaster.SendFunction, error) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	result := awset.result

	if awset.broadcast == nil {
		return nil, NoBroadcastHandler
	}

	handler := awset.broadcast
	awset.result = nil

	sendFunction := func() error {
		err := handler.Broadcast(replicaID, name, result)
		if err != nil {
			awset.lock.Lock()
			defer awset.lock.Unlock()

			if awset.result != nil {
				result.Join(awset.result)
			}
			awset.result = result

			return err
		}

		return nil
	}

	return sendFunction, nil
}

func (awset *Aworset) Update(data interface{}) (broadcaster.UpdateFunction, error) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

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
func (awset *Aworset) Join(interface{}) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.Join(awset)
}
