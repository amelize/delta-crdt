package ccounter

import (
	"testing"
)

func TestCCounter_Inc(t *testing.T) {
	a := NewIntCounter(1)
	b := NewIntCounter(2)

	b.Join(a.Inc(2))
	a.Join(b.Inc(3))

	if a.Value() != 5 {
		t.Fatalf("expected %d but has %d", 5, a.Value())
	}
}
