package crdt

import (
	"sync"

	"github.com/delta-crdt/aworset"
)

type Aworset struct {
	set    *aworset.AWORSet
	result *aworset.AWORSet
	lock   sync.RWMutex
}

func NewAworset(replica string) *Aworset {
	return &Aworset{
		set: aworset.New(replica),
	}
}

func (awset *Aworset) Add(val interface{}) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	change := awset.set.Add(val)

	if awset.result != nil {
		awset.result.Join(change)
	} else {
		awset.result = change
	}
}

func (awset *Aworset) Remove(val interface{}) {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	change := awset.set.Remove(val)

	if awset.result != nil {
		awset.result.Join(change)
	} else {
		awset.result = change
	}
}

func (awset *Aworset) Reset() {
	awset.lock.Lock()
	defer awset.lock.Unlock()

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
func (awset *Aworset) GetChanges() interface{} {
	awset.lock.Lock()
	defer awset.lock.Unlock()

	result := awset.result
	awset.result = nil

	return result
}

// Join joins broadcasted changes into set
func (awset *Aworset) Join(interface{}) {
	awset.Join(awset)
}
