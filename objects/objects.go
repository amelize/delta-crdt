package objects

import (
	"errors"

	"github.com/delta-crdt/encoding"
	"github.com/delta-crdt/kernel"
)

var NotExists = errors.New("Not exists")

type Record struct {
	crdt    kernel.Joinable
	encoder encoding.Encoder
}

type Objects struct {
	objects map[string]Record
}

func (objs *Objects) Add(name string, crdt kernel.Joinable, encoder encoding.Encoder) {
	objs.objects[name] = Record{crdt: crdt, encoder: encoder}
}

func (objs *Objects) Get(name string) kernel.Joinable {
	record := objs.GetRecord(name)
	if record != nil {
		return record.crdt
	}

	return nil
}

func (objs *Objects) GetRecord(name string) *Record {
	record, exists := objs.objects[name]
	if exists {
		return &record
	}

	return nil
}

func (objs *Objects) GetEncoded(name string) ([]byte, error) {
	record := objs.GetRecord(name)

	if record != nil {

		// bytes, err := record.encoder.Encode()
	}

	return nil, NotExists
}
