package crdt

import (
	"testing"
)

func TestReplica_CreateNewAWORSet(t *testing.T) {
	replicaOne := NewReplica("a")
	replicaTwo := NewReplica("b")

	setOne := replicaOne.CreateNewAWORSet("user.set", nil)
	setTwo := replicaTwo.CreateNewAWORSet("user.set", nil)

	setOne.Add("HelloBadge")
	setTwo.Add("WelocomeBadge")

	
}
