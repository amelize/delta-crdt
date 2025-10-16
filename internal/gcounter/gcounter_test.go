package gcounter

import "testing"

func TestGCounter_Simple(t *testing.T) {
	counterOne := New("c-1")
	counterTwo := New("c-2")
	counterTree := New("c-3")

	st1 := counterOne.Inc(1)
	st2 := counterTwo.Inc(1)
	st3 := counterTree.Inc(3)

	counterOne.Join(st2)
	counterOne.Join(st3)

	counterTwo.Join(st1)
	counterTwo.Join(st3)

	counterTree.Join(st1)
	counterTree.Join(st2)

	if counterOne.Value() != 5 {
		t.Fatalf("expect 5 but have %d", counterOne.Value())
	}

	if counterTwo.Value() != 5 {
		t.Fatalf("expect 5 but have %d", counterTwo.Value())
	}

	if counterTree.Value() != 5 {
		t.Fatalf("expect 5 but have %d", counterTree.Value())
	}
}
