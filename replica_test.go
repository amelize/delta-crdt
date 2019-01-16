package crdt

import (
	"encoding/json"
	"testing"

	"github.com/delta-crdt/kernel"

	"github.com/delta-crdt/aworset"
)

type KernelData struct {
	Pair  kernel.Pair
	Value string
}

type JsonData struct {
	Context kernel.ContextData
	Data    []KernelData
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

func TestReplica_CreateNewAWORSet(t *testing.T) {
	lock := make(chan struct{})

	replicaOne := NewReplica("a")
	replicaTwo := NewReplica("b")

	firstHandler := DummyHandler{other: replicaTwo}
	secondHandler := DummyHandler{other: replicaOne}

	setOne := replicaOne.CreateNewAWORSet("user.set", firstHandler)
	setTwo := replicaTwo.CreateNewAWORSet("user.set", secondHandler)

	setTwo.SetOnUpdated(func() {
		lock <- struct{}{}
	})

	setOne.Add("HelloBadge")
	setTwo.Add("WelocomeBadge")

	<-lock
}
