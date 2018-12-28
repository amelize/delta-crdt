package aworset

import (
	"log"
	"testing"
)

func TestAWORSet_Add(t *testing.T) {
	set1 := New("set-1")
	set2 := New("set-2")

	res1 := set1.Add("hello")
	res2 := set2.Add("world")

	set1.Join(res1)
	set1.Join(res2)

	set2.Join(set1)
	set2.Join(set2)

	setRes1 := set1.Value()
	setRes2 := set2.Value()

	log.Printf("res1: %#+v", setRes1)
	log.Printf("res2: %#+v", setRes2)

	_, ok := setRes1["hello"]
	if !ok {
		t.Fail()
	}

	_, ok = setRes1["world"]
	if !ok {
		t.Fail()
	}

	_, ok = setRes2["hello"]
	if !ok {
		t.Fail()
	}

	_, ok = setRes2["world"]
	if !ok {
		t.Fail()
	}

	res1a := set1.Remove("hello")

	set2.Join(res1a)
	set2.Join(res1a)

	setRes1 = set1.Value()
	setRes2 = set2.Value()
	log.Printf("res1: %#+v", setRes1)
	log.Printf("res2: %#+v", setRes2)

	_, ok = setRes1["hello"]
	if ok {
		t.Fail()
	}

	_, ok = setRes1["world"]
	if !ok {
		t.Fail()
	}

	_, ok = setRes2["hello"]
	if ok {
		t.Fail()
	}

	_, ok = setRes2["world"]
	if !ok {
		t.Fail()
	}

}
