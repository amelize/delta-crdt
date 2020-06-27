package crdt

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/amelize/delta-crdt/ccounter"
	"github.com/amelize/delta-crdt/kernel"

	"github.com/amelize/delta-crdt/aworset"
)

type KernelData struct {
	Pair     kernel.Pair
	Value    string
	IntValue int64
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

	log.Printf("%s", string(bts))

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
	currentKernel := counter.GetCounter().GetKernel()

	log.Printf("value current %d", counter.Value())

	data := JsonData{
		Context: counter.Context().GetData(),
		Data:    make([]KernelData, 0),
		ID:      counter.GetId(),
	}

	it := currentKernel.Dots.GetIterator()

	for it.HasMore() {
		kData := KernelData{
			Pair:     it.Key().(kernel.Pair),
			IntValue: int64(it.Value().(ccounter.IntValue)),
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

func (handler DummyIntHandler) OnUpdate(data interface{}) (*ccounter.IntCounter, error) {
	jsonData := JsonData{}
	json.Unmarshal(data.([]byte), &jsonData)

	newKernel := kernel.NewDotKernel()
	newKernel.Ctx = kernel.NewFromData(jsonData.Context)

	for _, v := range jsonData.Data {
		newKernel.Dots.Insert(v.Pair, ccounter.IntValue(v.IntValue))
	}

	cnt := ccounter.NewIntCounterWithKernel(jsonData.ID, newKernel)

	log.Printf("-> %d", cnt.Value())

	return cnt, nil
}

func TestReplica_CreateNewAWORSet(t *testing.T) {
	lock := make(chan struct{})

	broadcastRate := time.Microsecond * 1500
	replicaOne := NewReplicaWithSelfUniqueAddress(broadcastRate)
	replicaTwo := NewReplicaWithSelfUniqueAddress(broadcastRate)

	firstHandler := DummyHandler{other: replicaTwo}
	secondHandler := DummyHandler{other: replicaOne}

	setOne := replicaOne.CreateNewAWORSet("user.set", firstHandler)
	setTwo := replicaTwo.CreateNewAWORSet("user.set", secondHandler)

	setTwo.SetOnUpdated(func() {
		lock <- struct{}{}
	})

	setOne.Add("Value-One")
	setOne.Add("Value-Two")
	setTwo.Add("R-One")
	setTwo.Add("R-Two")

	<-lock

	time.Sleep(1 * time.Second)

	log.Printf(": %+v", setOne.Value())
	log.Printf(": %+v", setTwo.Value())

	for k := range setOne.Value() {
		if !setTwo.In(k) {
			t.Errorf("Not found %s", k)
		}
	}

	for k := range setTwo.Value() {
		if !setOne.In(k) {
			t.Errorf("Not found %s", k)
		}
	}

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

	time.Sleep(1 * time.Second)

	if setOne.Value() != -5 {
		t.Fatalf("Wrong value %d", setOne.Value())
	}
}
