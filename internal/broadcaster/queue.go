package broadcaster

import (
	"container/list"
	"sync"
)

func newQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

type Queue struct {
	changeLock sync.Mutex
	list       *list.List
}

func (queue *Queue) Head() *string {
	queue.changeLock.Lock()
	defer queue.changeLock.Unlock()

	item := queue.list.Front()
	if item == nil {
		return nil
	}

	queue.list.Remove(item)
	itemString := item.Value.(string)

	return &itemString
}

// Adding change to queue
func (queue *Queue) Push(name string) {
	queue.changeLock.Lock()
	defer queue.changeLock.Unlock()

	queue.list.PushBack(name)
}
