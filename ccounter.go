package crdt

import (
	"sync"

	"github.com/amelize/delta-crdt/broadcaster"
	"github.com/amelize/delta-crdt/ccounter"
)

type CCounterBroadcastHandler interface {
	Broadcast(replicaID, name string, counter *ccounter.IntCounter) error
	OnUpdate(data interface{}) (*ccounter.IntCounter, error)
}

type CCounter struct {
	counter *ccounter.IntCounter
	result  *ccounter.IntCounter

	lock      sync.RWMutex
	broadcast CCounterBroadcastHandler

	onChanged broadcaster.OnChanged
	onUpdated broadcaster.OnUpdated
}

func NewCCounter(replica string, broadcastHandler CCounterBroadcastHandler) *CCounter {
	return &CCounter{
		counter:   ccounter.NewIntCounter(replica),
		broadcast: broadcastHandler,
	}
}

func (awset *CCounter) SetOnChanged(onChanged broadcaster.OnChanged) {
	awset.onChanged = onChanged
}

func (awset *CCounter) SetOnUpdated(onUpdated broadcaster.OnUpdated) {
	awset.onUpdated = onUpdated
}

func (awset *CCounter) SetBroadcastHandler(broadcast CCounterBroadcastHandler) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	awset.broadcast = broadcast
}

func (cnt *CCounter) Inc(val int64) {
	cnt.lock.Lock()
	defer func() {
		cnt.lock.Unlock()

		if cnt.onChanged != nil {
			go cnt.onChanged()
		}
	}()

	change := cnt.counter.Inc(val)

	if cnt.result != nil {
		cnt.result.Join(change)
	} else {
		cnt.result = change
	}
}

func (cnt *CCounter) Dec(val int64) {
	cnt.lock.Lock()
	defer func() {
		cnt.lock.Unlock()

		if cnt.onChanged != nil {
			go cnt.onChanged()
		}
	}()

	change := cnt.counter.Dec(val)

	if cnt.result != nil {
		cnt.result.Join(change)
	} else {
		cnt.result = change
	}
}

func (cnt *CCounter) Reset() {
	cnt.lock.Lock()
	defer func() {
		cnt.lock.Unlock()

		if cnt.onChanged != nil {
			go cnt.onChanged()
		}
	}()

	change := cnt.counter.Reset().(interface{})

	if cnt.result != nil {
		cnt.result.Join(change)
	} else {
		cnt.result = change.(*ccounter.IntCounter)
	}
}

func (cnt *CCounter) Value() int64 {
	cnt.lock.RLock()
	defer cnt.lock.RUnlock()

	return cnt.counter.Value()
}

// GetChanges returns changes for broadcast and clears changes.
func (cnt *CCounter) Broadcast(replicaID, name string) (broadcaster.SendFunction, error) {
	cnt.lock.RLock()
	defer func() {
		cnt.lock.RUnlock()
	}()

	result := cnt.result

	if cnt.broadcast == nil {
		return nil, NoBroadcastHandler
	}

	handler := cnt.broadcast
	cnt.result = nil

	sendFunction := func() error {
		err := handler.Broadcast(replicaID, name, result)
		if err != nil {
			cnt.lock.Lock()
			defer cnt.lock.Unlock()

			if cnt.result != nil {
				result.Join(cnt.result)
			}

			cnt.result = result

			return err
		}

		return nil
	}

	return sendFunction, nil
}

func (cnt *CCounter) Update(data interface{}) (broadcaster.UpdateFunction, error) {
	cnt.lock.RLock()
	defer func() {
		cnt.lock.RUnlock()
	}()

	if cnt.broadcast == nil {
		return nil, NoBroadcastHandler
	}

	updateFunction := func() error {
		set, err := cnt.broadcast.OnUpdate(data)
		if err != nil {
			return err
		}

		cnt.Join(set)

		if cnt.onUpdated != nil {
			cnt.onUpdated()
		}

		return nil
	}

	return updateFunction, nil
}

// Join joins broadcasted changes into set
func (cnt *CCounter) Join(data interface{}) {
	cnt.lock.Lock()
	defer cnt.lock.Unlock()

	cnt.counter.Join(data)
}
