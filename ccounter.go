package crdt

import (
	"sync"

	"github.com/amelize/delta-crdt/internal/broadcaster"
	"github.com/amelize/delta-crdt/internal/ccounter"
)

type CCounterBroadcastHandler interface {
	Broadcast(replicaID int64, name string, counter *ccounter.IntCounter) error
	OnUpdate(data interface{}) (*ccounter.IntCounter, error)
}

type CCounter struct {
	counter *ccounter.IntCounter
	result  *ccounter.IntCounter

	lock      sync.RWMutex
	broadcast CCounterBroadcastHandler
	name      string

	onChanged broadcaster.ChangeHandlerInterface
	onUpdated broadcaster.OnUpdated
}

func NewCCounter(replica int64, name string, broadcastHandler CCounterBroadcastHandler) *CCounter {
	return &CCounter{
		counter:   ccounter.NewIntCounter(replica),
		broadcast: broadcastHandler,
		name:      name,
	}
}

func (awset *CCounter) GetName() string {
	return awset.name
}

func (awset *CCounter) SetOnChanged(onChanged broadcaster.ChangeHandlerInterface) {
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
			cnt.onChanged.OnChange(cnt.name)
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
			cnt.onChanged.OnChange(cnt.name)
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
			cnt.onChanged.OnChange(cnt.name)
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
func (cnt *CCounter) Broadcast(replicaID int64, name string) error {
	cnt.lock.RLock()
	defer func() {
		cnt.lock.RUnlock()
	}()

	result := cnt.result

	if cnt.broadcast == nil {
		return NoBroadcastHandler
	}

	handler := cnt.broadcast
	cnt.result = nil

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

func (cnt *CCounter) Update(data interface{}) error {
	cnt.lock.RLock()
	defer func() {
		cnt.lock.RUnlock()
	}()

	if cnt.broadcast == nil {
		return NoBroadcastHandler
	}

	set, err := cnt.broadcast.OnUpdate(data)
	if err != nil {
		return err
	}

	cnt.counter.Join(set)

	if cnt.onUpdated != nil {
		go cnt.onUpdated()
	}

	return nil
}

// Join joins broadcasted changes into set
func (cnt *CCounter) Join(data interface{}) {
	cnt.lock.Lock()
	defer cnt.lock.Unlock()

	cnt.counter.Join(data)
}
