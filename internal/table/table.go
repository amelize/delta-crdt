package table

import (
	"sync"

	"github.com/amelize/delta-crdt/internal/object"
)

type Table struct {

	// Keeps objects
	objects map[uint64]*object.Object
	// Map name to ID
	objectNames map[string]uint64

	// General lock to change objects table
	lock sync.RWMutex
}

func NewTable() Table {
	return Table{
		objects:     make(map[uint64]*object.Object, 10000),
		objectNames: make(map[string]uint64, 10000),
	}
}

func (table *Table) Get(name string) *object.Object {
	table.lock.RLock()
	defer table.lock.RUnlock()

	id, exists := table.objectNames[name]
	if !exists {
		return nil
	}

	return table.objects[id]
}

func (table *Table) Set(name string, obj object.Object) *object.Object {
	table.lock.Lock()
	defer table.lock.Unlock()

	id, exists := table.objectNames[name]
	if !exists {
		return table.createNewEntry(name, obj)
	}

	return table.updateData(id, obj)
}

func (table *Table) createNewId() uint64 {

}

func (table *Table) updateData(id uint64, obj object.Object) *object.Object {
	table.objects[id].InnterType.Update(obj.InnterType)
}

func (table *Table) createNewEntry(name string, obj object.Object) *object.Object {
	id := table.createNewId()

	obj.ObjectID = id

	table.objectNames[name] = id
	table.objects[id] = &obj

	return &obj
}
