package broadcaster

import (
	"errors"
	"sync"
)

var NotExists = errors.New("Not exists")

type SendFunction = func() error
type UpdateFunction = func() error
type OnChanged = func()
type OnUpdated = func()

type Broadcastable interface {
	Broadcast(replicaID, name string) (SendFunction, error)
	Update(data interface{}) (UpdateFunction, error)
	SetOnChanged(onChanged OnChanged)
}

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
	crdt.SetOnChanged(func() {
		objs.change.Push(name)
	})
}

func (objs *Objects) Get(name string) Broadcastable {
	record, exists := objs.objects[name]
	if exists {
		return record
	}

	return nil
}

func (objs *Objects) GetChangedHead() string {
	return objs.change.Head()
}

func (objs *Objects) Resend(name string) {
	objs.change.Push(name)
}
