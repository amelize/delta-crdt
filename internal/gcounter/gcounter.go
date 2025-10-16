package gcounter

// GCounter Simple counter
type GCounter struct {
	state     map[string]int64
	replicaID string
}

// New Create new counter
func New(id string) *GCounter {
	return &GCounter{
		replicaID: id,
		state:     make(map[string]int64),
	}
}

// New Create new counter
func empty() *GCounter {
	return &GCounter{
		state: make(map[string]int64),
	}
}

func (counter GCounter) Local() int64 {
	return counter.state[counter.replicaID]
}

func (counter GCounter) Value() int64 {
	var result int64

	for _, v := range counter.state {
		result += v
	}

	return result
}

// Inc Increment value by value
func (counter GCounter) Inc(value int64) *GCounter {
	changedValue := counter.state[counter.replicaID] + value
	counter.state[counter.replicaID] = changedValue

	change := empty()
	change.state[counter.replicaID] = changedValue

	return change
}

func (counter GCounter) Join(other interface{}) {
	otherCounter, ok := other.(*GCounter)

	if ok {
		for k, v := range otherCounter.state {
			current := counter.state[k]
			if current < v {
				counter.state[k] = v
			}
		}
	} else {
		panic("incorrect join")
	}
}
