package broadcaster

import (
	"container/list"
	"sync"
)

type element struct {
	next *element
}

func newQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

type Queue struct {
	changeLock sync.Mutex
	list       *list.List
}

func (queue *Queue) Head() string {
	queue.changeLock.Lock()
	defer queue.changeLock.Unlock()

	item := queue.list.Front()
	if item != nil {
		queue.list.Remove(item)
	} else {
		return ""
	}

	return item.Value.(string)
}

func (queue *Queue) Push(name string) {
	queue.changeLock.Lock()
	defer queue.changeLock.Unlock()

	queue.list.PushBack(name)
}
