package kernel

import (
	"testing"
)

func TestDotContext_makeDot(t *testing.T) {
	ctx := NewDotContext()

	dotOne := ctx.makeDot(1)
	dotTwo := ctx.makeDot(1)
	dotThree := ctx.makeDot(1)

	if dotOne.Second == dotTwo.Second {
		t.Fatalf("Same value for dots")
	}

	if dotOne.Second == dotThree.Second {
		t.Fatalf("Same value for dots")
	}
}
