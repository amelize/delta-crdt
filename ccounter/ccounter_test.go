package ccounter

import (
	"testing"
)

func TestCCounter_Inc(t *testing.T) {
	a := NewIntCounter("a")
	b := NewIntCounter("b")

	b.Join(a.Inc(2))
	a.Join(b.Inc(3))

	if a.Value() != 5 {
		t.Fatalf("expected %d but has %d", 5, a.Value())
	}
}
