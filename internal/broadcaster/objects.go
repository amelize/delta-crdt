package broadcaster

import (
	"sync"
)

// Map of objects
type Objects struct {
	objects map[string]Broadcastable
	lock    sync.RWMutex

	change *Queue
}

func NewObjects() *Objects {
	return &Objects{
		objects: make(map[string]Broadcastable),
		change:  newQueue(),
	}
}

func (objs *Objects) Add(name string, crdt Broadcastable) {
	objs.lock.Lock()
	defer objs.lock.Unlock()

	objs.objects[name] = crdt
	crdt.SetOnChanged(objs)
}

func (objs *Objects) Get(name string) Broadcastable {
	objs.lock.RLock()
	defer objs.lock.RUnlock()

	record, exists := objs.objects[name]
	if exists {
		return record
	}

	return nil
}

func (objs *Objects) OnChange(name string) {
	// Adding change
	objs.change.Push(name)
}

// Get first changed item
func (objs *Objects) GetChangedHead() *string {
	objs.lock.Lock()
	defer objs.lock.Unlock()

	return objs.change.Head()
}

func (objs *Objects) Resend(name string) {
	objs.lock.Lock()
	defer objs.lock.Unlock()

	objs.change.Push(name)
}
