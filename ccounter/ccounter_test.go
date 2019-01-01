package ccounter

import (
	"testing"
)

func TestCCounter_Inc(t *testing.T) {
	a := NewIntCounter("a")
	b := NewIntCounter("b")

	b.Join(a.Inc(1))
	a.Join(b.Inc(1))

	if a.Value() != 2 {
		t.Fatalf("expected %d but has %d", 2, a.Value())
	}
}
