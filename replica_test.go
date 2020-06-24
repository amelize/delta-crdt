package crdt

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/amelize/delta-crdt/ccounter"
	"github.com/amelize/delta-crdt/kernel"

	"github.com/amelize/delta-crdt/aworset"
)

type KernelData struct {
	Pair  kernel.Pair
	Value string
}

type JsonData struct {
	Context kernel.ContextData
	Data    []KernelData
	ID      string
}

type DummyHandler struct {
	other *Replica
}

func (handler DummyHandler) Broadcast(replicaID, name string, aworset *aworset.AWORSet) error {
	currentKernel := aworset.GetKernel()

	data := JsonData{
		Context: currentKernel.Ctx.GetData(),
		Data:    make([]KernelData, 0),
	}

	it := currentKernel.Dots.GetIterator()

	for it.HasMore() {
		kData := KernelData{
			Pair:  it.Key().(kernel.Pair),
			Value: it.Value().(string),
		}

		data.Data = append(data.Data, kData)
		it.Next()
	}

	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return handler.other.Update(name, bts)
}

func (handler DummyHandler) OnUpdate(data interface{}) (*aworset.AWORSet, error) {
	jsonData := JsonData{}
	json.Unmarshal(data.([]byte), &jsonData)

	newKernel := kernel.NewDotKernel()
	newKernel.Ctx = kernel.NewFromData(jsonData.Context)

	for _, v := range jsonData.Data {
		newKernel.Dots.Insert(v.Pair, v.Value)
	}

	set := aworset.NewFromKernel(newKernel)

	return set, nil
}

type DummyIntHandler struct {
	other *Replica
}

func (handler DummyIntHandler) Broadcast(replicaID, name string, counter *ccounter.IntCounter) error {
	currentKernel := counter.Context()

	data := JsonData{
		Context: currentKernel.GetData(),
		Data:    make([]KernelData, 0),
		ID:      counter.GetId(),
	}

	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return handler.other.Update(name, bts)
}

func (handler DummyIntHandler) OnUpdate(data interface{}) (*ccounter.IntCounter, error) {
	jsonData := JsonData{}
	json.Unmarshal(data.([]byte), &jsonData)

	newKernel := kernel.NewDotKernel()
	newKernel.Ctx = kernel.NewFromData(jsonData.Context)

	cnt := ccounter.NewIntCounterWithContex(jsonData.ID, newKernel.Ctx)

	return cnt, nil
}

func TestReplica_CreateNewAWORSet(t *testing.T) {
	lock := make(chan struct{})

	broadcastRate := time.Millisecond * 500
	replicaOne := NewReplica("a", broadcastRate)
	replicaTwo := NewReplica("b", broadcastRate)

	firstHandler := DummyHandler{other: replicaTwo}
	secondHandler := DummyHandler{other: replicaOne}

	setOne := replicaOne.CreateNewAWORSet("user.set", firstHandler)
	setTwo := replicaTwo.CreateNewAWORSet("user.set", secondHandler)

	setTwo.SetOnUpdated(func() {
		lock <- struct{}{}
	})

	setOne.Add("HelloBadge")
	setTwo.Add("WelocomeBadge")
	// setTwo.Add("One more")

	<-lock
}

func TestReplica_CreateCCounter(t *testing.T) {
	lock := make(chan struct{})

	broadcastRate := time.Millisecond * 500
	replicaOne := NewReplica("a", broadcastRate)
	replicaTwo := NewReplica("b", broadcastRate)

	firstHandler := DummyIntHandler{other: replicaTwo}
	secondHandler := DummyIntHandler{other: replicaOne}

	setOne := replicaOne.CreateCCounter("user.setX", firstHandler)
	setTwo := replicaTwo.CreateCCounter("user.setX", secondHandler)

	setTwo.SetOnUpdated(func() {
		lock <- struct{}{}
	})

	setOne.Inc(10)
	setTwo.Dec(15)

	<-lock

	if setOne.Value() != -5 {
		t.Fatalf("Wrong value %d", setOne.Value())
	}
}
